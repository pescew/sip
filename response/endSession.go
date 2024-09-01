package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDEndSession = "36"

var ErrInvalidResponse36 = fmt.Errorf("Invalid SIP %s response", MsgIDEndSession)

// The ACS must send this message in response to the End Patron Session message.
type EndSession struct {
	// Required Fields:
	EndSession      bool
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional Fields:
	ScreenMessage string `validate:"sip"`
	PrintLine     string `validate:"sip"`
}

func (es *EndSession) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDEndSession)
	msg.WriteString(utils.YorN(es.EndSession))
	msg.WriteString(es.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", es.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", es.PatronID, delimiter)

	if es.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", es.ScreenMessage, delimiter)
	}

	if es.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", es.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (es *EndSession) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 27 {
		return 0, ErrInvalidResponse36
	}

	if string(runes[0:2]) != MsgIDEndSession {
		return 0, ErrInvalidResponse36
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	es.EndSession = utils.ParseBool(runes[2])

	es.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	es.InstitutionID = codes["AO"]
	es.PatronID = codes["AA"]

	es.ScreenMessage = codes["AF"]
	es.PrintLine = codes["AG"]

	err = es.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (es *EndSession) Validate() error {
	err := Validate.Struct(es)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDEndSession, err.(validator.ValidationErrors))
	}
	return nil
}

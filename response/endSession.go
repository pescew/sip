package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse36 = fmt.Errorf("Invalid SIP %s", types.RespEndSession.String())

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

	SeqNum int `validate:"min=0,max=9"`
}

func (es *EndSession) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespEndSession.ID())
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

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", es.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (es *EndSession) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 27 {
		return ErrInvalidResponse36
	}

	if string(runes[0:2]) != types.RespEndSession.ID() {
		return ErrInvalidResponse36
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		es.SeqNum = 0
	} else {
		es.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			es.SeqNum = 0
		}
	}

	es.EndSession = utils.ParseBool(runes[2])

	es.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return err
	}

	es.InstitutionID = codes["AO"]
	es.PatronID = codes["AA"]

	es.ScreenMessage = codes["AF"]
	es.PrintLine = codes["AG"]

	err = es.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (es *EndSession) Validate() error {
	err := Validate.Struct(es)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespEndSession.String(), err.(validator.ValidationErrors))
	}
	return nil
}

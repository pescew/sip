package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

const MsgIDPatronEnable = "26"

var ErrInvalidResponse26 = fmt.Errorf("Invalid SIP %s response", MsgIDPatronEnable)

// The ACS must send this message in response to the Fee Paid message.
type PatronEnable struct {
	// Required Fields:
	PatronStatus    fields.PatronStatus `validate:"required"`
	Language        int                 `validate:"min=0,max=999"`
	TransactionDate time.Time           `validate:"required"`
	InstitutionID   string              `validate:"sip"`
	PatronID        string              `validate:"sip"`
	PatronName      string              `validate:"sip"`

	// Optional Fields:
	ValidPatron         bool
	ValidPatronPassword bool
	ScreenMessage       string `validate:"sip"`
	PrintLine           string `validate:"sip"`
}

func (pe *PatronEnable) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDPatronEnable)

	msg.WriteString(pe.PatronStatus.Marshal())
	fmt.Fprintf(&msg, "%03d", pe.Language)
	msg.WriteString(pe.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", pe.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", pe.PatronID, delimiter)
	fmt.Fprintf(&msg, "AE%s%c", pe.PatronName, delimiter)

	fmt.Fprintf(&msg, "BL%s%c", utils.YorN(pe.ValidPatron), delimiter)
	fmt.Fprintf(&msg, "CQ%s%c", utils.YorN(pe.ValidPatronPassword), delimiter)

	if pe.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", pe.ScreenMessage, delimiter)
	}

	if pe.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", pe.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (pe *PatronEnable) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 46 {
		return 0, ErrInvalidResponse26
	}

	if string(runes[0:2]) != MsgIDPatronEnable {
		return 0, ErrInvalidResponse26
	}

	codes := utils.ExtractFields(string(runes[37:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AE": "", "BL": "", "CQ": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	err = pe.PatronStatus.Unmarshal(string(runes[2:16]))
	if err != nil {
		return 0, err
	}

	pe.Language, err = strconv.Atoi(string(runes[16:19]))
	if err != nil {
		return 0, err
	}

	pe.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[19:37]))
	if err != nil {
		return 0, err
	}

	pe.InstitutionID = codes["AO"]
	pe.PatronID = codes["AA"]
	pe.PatronName = codes["AE"]

	if codes["BL"] != "" {
		pe.ValidPatron = utils.ParseBool([]rune(codes["BL"])[0])
	}

	if codes["CQ"] != "" {
		pe.ValidPatronPassword = utils.ParseBool([]rune(codes["CQ"])[0])
	}

	pe.ScreenMessage = codes["AF"]
	pe.PrintLine = codes["AG"]

	err = pe.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (pe *PatronEnable) Validate() error {
	err := Validate.Struct(pe)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDPatronEnable, err.(validator.ValidationErrors))
	}
	return nil
}

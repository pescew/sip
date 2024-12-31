package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest23 = fmt.Errorf("Invalid SIP %s request", fields.ReqPatronStatus.String())

// This message is used by the SC to request patron information from the ACS. The ACS must respond to this command with a Patron Status Response message.
type PatronStatus struct {
	// Required:
	Language         fields.LanguageCode `validate:"min=0,max=999"`
	TransactionDate  time.Time           `validate:"required"`
	InstitutionID    string              `validate:"sip"`
	PatronID         string              `validate:"required,sip"`
	TerminalPassword string              `validate:"sip"`
	PatronPassword   string              `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (ps *PatronStatus) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(fields.ReqPatronStatus.ID())

	msg.WriteString(ps.Language.ID())
	msg.WriteString(ps.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ps.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", ps.PatronID, delimiter)

	fmt.Fprintf(&msg, "AC%s%c", ps.TerminalPassword, delimiter)
	fmt.Fprintf(&msg, "AD%s%c", ps.PatronPassword, delimiter)

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", ps.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (ps *PatronStatus) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 35 {
		return ErrInvalidRequest23
	}

	if string(runes[0:2]) != fields.ReqPatronStatus.ID() {
		return ErrInvalidRequest23
	}

	codes := utils.ExtractFields(string(runes[23:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ps.SeqNum = 0
	} else {
		ps.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ps.SeqNum = 0
		}
	}

	language, err := strconv.Atoi(string(runes[2:5]))
	if err != nil {
		return err
	}
	ps.Language = fields.LanguageCode(language)

	ps.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[5:23]))
	if err != nil {
		return err
	}

	ps.InstitutionID = codes["AO"]
	ps.PatronID = codes["AA"]
	ps.TerminalPassword = codes["AC"]
	ps.PatronPassword = codes["AD"]

	err = ps.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (ps *PatronStatus) Validate() error {
	err := Validate.Struct(ps)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", fields.ReqPatronStatus.String(), err.(validator.ValidationErrors))
	}
	return nil
}

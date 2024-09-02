package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest25 = fmt.Errorf("Invalid SIP %s request", types.ReqPatronEnable.String())

// This message can be used by the SC to re-enable canceled patrons. It should only be used for system testing and validation. The ACS should respond with a Patron Enable Response message.
type PatronEnable struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (pe *PatronEnable) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqPatronEnable.ID())

	msg.WriteString(pe.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", pe.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", pe.PatronID, delimiter)

	if pe.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", pe.TerminalPassword, delimiter)
	}

	if pe.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", pe.PatronPassword, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", pe.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (pe *PatronEnable) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 52 {
		return ErrInvalidRequest25
	}

	if string(runes[0:2]) != types.ReqPatronEnable.ID() {
		return ErrInvalidRequest25
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		pe.SeqNum = 0
	} else {
		pe.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			pe.SeqNum = 0
		}
	}

	pe.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	pe.InstitutionID = codes["AO"]
	pe.PatronID = codes["AA"]
	pe.TerminalPassword = codes["AC"]
	pe.PatronPassword = codes["AD"]

	pe.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (pe *PatronEnable) Validate() error {
	err := Validate.Struct(pe)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqPatronEnable.String(), err.(validator.ValidationErrors))
	}
	return nil
}

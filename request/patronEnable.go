package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDPatronEnable = "25"

var ErrInvalidRequest25 = fmt.Errorf("Invalid SIP %s request", MsgIDPatronEnable)

// This message can be used by the SC to re-enable canceled patrons. It should only be used for system testing and validation. The ACS should respond with a Patron Enable Response message.
type PatronEnable struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`
}

func (pe *PatronEnable) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDPatronEnable)

	msg.WriteString(pe.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", pe.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", pe.PatronID, delimiter)

	if pe.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", pe.TerminalPassword, delimiter)
	}

	if pe.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", pe.PatronPassword, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (pe *PatronEnable) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 52 {
		return 0, ErrInvalidRequest25
	}

	if string(runes[0:2]) != MsgIDPatronEnable {
		return 0, ErrInvalidRequest25
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	pe.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
	}

	pe.InstitutionID = codes["AO"]
	pe.PatronID = codes["AA"]
	pe.TerminalPassword = codes["AC"]
	pe.PatronPassword = codes["AD"]

	pe.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (pe *PatronEnable) Validate() error {
	err := Validate.Struct(pe)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDPatronEnable, err.(validator.ValidationErrors))
	}
	return nil
}

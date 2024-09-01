package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDPatronStatus = "23"

var ErrInvalidRequest23 = fmt.Errorf("Invalid SIP %s request", MsgIDPatronStatus)

// This message is used by the SC to request patron information from the ACS. The ACS must respond to this command with a Patron Status Response message.
type PatronStatus struct {
	// Required:
	Language         int       `validate:"min=0,max=999"`
	TransactionDate  time.Time `validate:"required"`
	InstitutionID    string    `validate:"required,sip"`
	PatronID         string    `validate:"required,sip"`
	TerminalPassword string    `validate:"sip"`
	PatronPassword   string    `validate:"sip"`
}

func (ps *PatronStatus) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDPatronStatus)

	fmt.Fprintf(&msg, "%03d", ps.Language)
	msg.WriteString(ps.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ps.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", ps.PatronID, delimiter)

	fmt.Fprintf(&msg, "AC%s%c", ps.TerminalPassword, delimiter)
	fmt.Fprintf(&msg, "AD%s%c", ps.PatronPassword, delimiter)

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ps *PatronStatus) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 35 {
		return 0, ErrInvalidRequest23
	}

	if string(runes[0:2]) != MsgIDPatronStatus {
		return 0, ErrInvalidRequest23
	}

	codes := utils.ExtractFields(string(runes[23:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	ps.Language, err = strconv.Atoi(string(runes[2:5]))
	if err != nil {
		return 0, err
	}

	ps.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[5:23]))
	if err != nil {
		return 0, err
	}

	ps.InstitutionID = codes["AO"]
	ps.PatronID = codes["AA"]
	ps.TerminalPassword = codes["AC"]
	ps.PatronPassword = codes["AD"]

	err = ps.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (ps *PatronStatus) Validate() error {
	err := Validate.Struct(ps)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDPatronStatus, err.(validator.ValidationErrors))
	}
	return nil
}

package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDEndPatronSession = "35"

var ErrInvalidRequest35 = fmt.Errorf("Invalid SIP %s request", MsgIDEndPatronSession)

// This message will be sent when a patron has completed all of their transactions. The ACS may, upon receipt of this command, close any open files or deallocate data structures pertaining to that patron. The ACS should respond with an End Session Response message.
type EndPatronSession struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`
}

func (eps *EndPatronSession) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDEndPatronSession)

	msg.WriteString(eps.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", eps.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", eps.PatronID, delimiter)

	if eps.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", eps.TerminalPassword, delimiter)
	}

	if eps.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", eps.PatronPassword, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (eps *EndPatronSession) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 52 {
		return 0, ErrInvalidRequest35
	}

	if string(runes[0:2]) != MsgIDEndPatronSession {
		return 0, ErrInvalidRequest35
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

	eps.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
	}

	eps.InstitutionID = codes["AO"]
	eps.PatronID = codes["AA"]
	eps.TerminalPassword = codes["AC"]
	eps.PatronPassword = codes["AD"]

	eps.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (eps *EndPatronSession) Validate() error {
	err := Validate.Struct(eps)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDEndPatronSession, err.(validator.ValidationErrors))
	}
	return nil
}

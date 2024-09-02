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

var ErrInvalidRequest35 = fmt.Errorf("Invalid SIP %s request", types.ReqEndPatronSession.String())

// This message will be sent when a patron has completed all of their transactions. The ACS may, upon receipt of this command, close any open files or deallocate data structures pertaining to that patron. The ACS should respond with an End Session Response message.
type EndPatronSession struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (eps *EndPatronSession) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqEndPatronSession.ID())

	msg.WriteString(eps.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", eps.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", eps.PatronID, delimiter)

	if eps.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", eps.TerminalPassword, delimiter)
	}

	if eps.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", eps.PatronPassword, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", eps.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (eps *EndPatronSession) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 52 {
		return ErrInvalidRequest35
	}

	if string(runes[0:2]) != types.ReqEndPatronSession.ID() {
		return ErrInvalidRequest35
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		eps.SeqNum = 0
	} else {
		eps.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			eps.SeqNum = 0
		}
	}

	eps.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	eps.InstitutionID = codes["AO"]
	eps.PatronID = codes["AA"]
	eps.TerminalPassword = codes["AC"]
	eps.PatronPassword = codes["AD"]

	eps.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (eps *EndPatronSession) Validate() error {
	err := Validate.Struct(eps)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqEndPatronSession.String(), err.(validator.ValidationErrors))
	}
	return nil
}

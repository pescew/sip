package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest65 = fmt.Errorf("Invalid SIP %s request", types.ReqRenewAll.String())

// This message is used to renew all items that the patron has checked out. The ACS should respond with a Renew All Response message.
type RenewAll struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	PatronPassword   string `validate:"sip"`
	TerminalPassword string `validate:"sip"`
	FeeAcknowledged  bool

	SeqNum int `validate:"min=0,max=9"`
}

func (ra *RenewAll) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqRenewAll.ID())

	msg.WriteString(ra.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ra.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", ra.PatronID, delimiter)

	if ra.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", ra.PatronPassword, delimiter)
	}

	if ra.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", ra.TerminalPassword, delimiter)
	}

	if ra.FeeAcknowledged {
		fmt.Fprintf(&msg, "BO%s%c", utils.YorN(ra.FeeAcknowledged), delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", ra.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (ra *RenewAll) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 26 {
		return ErrInvalidRequest65
	}

	if string(runes[0:2]) != types.ReqRenewAll.ID() {
		return ErrInvalidRequest65
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AD": "", "AC": "", "BO": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ra.SeqNum = 0
	} else {
		ra.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ra.SeqNum = 0
		}
	}

	ra.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	ra.InstitutionID = codes["AO"]
	ra.PatronID = codes["AA"]
	ra.PatronPassword = codes["AD"]
	ra.TerminalPassword = codes["AC"]

	if utf8.RuneCountInString(codes["BO"]) > 0 {
		ra.FeeAcknowledged = utils.ParseBool([]rune(codes["BO"])[0])
	}

	err = ra.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (ra *RenewAll) Validate() error {
	err := Validate.Struct(ra)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqRenewAll.String(), err.(validator.ValidationErrors))
	}
	return nil
}

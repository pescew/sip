package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDRenewAll = "65"

var ErrInvalidRequest65 = fmt.Errorf("Invalid SIP %s request", MsgIDRenewAll)

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
}

func (ra *RenewAll) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDRenewAll)

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

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ra *RenewAll) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 26 {
		return 0, ErrInvalidRequest65
	}

	if string(runes[0:2]) != MsgIDRenewAll {
		return 0, ErrInvalidRequest65
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AD": "", "AC": "", "BO": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	ra.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
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
		return 0, err
	}

	return seqNum, nil
}

func (ra *RenewAll) Validate() error {
	err := Validate.Struct(ra)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDRenewAll, err.(validator.ValidationErrors))
	}
	return nil
}

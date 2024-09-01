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

const MsgIDCheckout = "11"

var ErrInvalidRequest11 = fmt.Errorf("Invalid SIP %s request", MsgIDCheckout)

// This message is used by the SC to request to check out an item, and also to cancel a Checkin request that did not successfully complete. The ACS must respond to this command with a Checkout Response message.
type Checkout struct {
	// Required:
	SCRenewalPolicy  bool
	NoBlock          bool
	TransactionDate  time.Time `validate:"required"`
	NBDueDate        time.Time `validate:"required"`
	InstitutionID    string    `validate:"required,sip"`
	PatronID         string    `validate:"required,sip"`
	ItemID           string    `validate:"required,sip"`
	TerminalPassword string    `validate:"sip"`

	// Optional:
	ItemProperties  string `validate:"sip"`
	PatronPassword  string `validate:"sip"`
	FeeAcknowledged bool
	Cancel          bool
}

func (co *Checkout) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDCheckout)

	msg.WriteString(utils.YorN(co.SCRenewalPolicy))
	msg.WriteString(utils.YorN(co.NoBlock))

	msg.WriteString(co.TransactionDate.Format(utils.SIPDateFormat))
	msg.WriteString(co.NBDueDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", co.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", co.PatronID, delimiter)

	fmt.Fprintf(&msg, "AB%s%c", co.ItemID, delimiter)
	fmt.Fprintf(&msg, "AC%s%c", co.TerminalPassword, delimiter)

	if co.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", co.ItemProperties, delimiter)
	}

	if co.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", co.PatronPassword, delimiter)
	}

	if co.FeeAcknowledged {
		fmt.Fprintf(&msg, "BO%s%c", utils.YorN(co.FeeAcknowledged), delimiter)
	}
	if co.Cancel {
		fmt.Fprintf(&msg, "BI%s%c", utils.YorN(co.Cancel), delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (co *Checkout) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 52 {
		return 0, ErrInvalidRequest11
	}

	if string(runes[0:2]) != MsgIDCheckout {
		return 0, ErrInvalidRequest11
	}

	codes := utils.ExtractFields(string(runes[40:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AB": "", "AC": "", "CH": "", "AD": "", "BO": "", "BI": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	co.SCRenewalPolicy = utils.ParseBool(runes[2])
	co.NoBlock = utils.ParseBool(runes[3])

	co.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[4:22]))
	if err != nil {
		return 0, err
	}

	co.NBDueDate, err = time.Parse(utils.SIPDateFormat, string(runes[22:40]))
	if err != nil {
		return 0, err
	}

	co.InstitutionID = codes["AO"]
	co.PatronID = codes["AA"]
	co.ItemID = codes["AB"]
	co.TerminalPassword = codes["AC"]

	co.ItemProperties = codes["CH"]
	co.PatronPassword = codes["AD"]

	if utf8.RuneCountInString(codes["BO"]) > 0 {
		co.FeeAcknowledged = utils.ParseBool([]rune(codes["BO"])[0])
	}

	if utf8.RuneCountInString(codes["BI"]) > 0 {
		co.Cancel = utils.ParseBool([]rune(codes["BI"])[0])
	}

	co.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (co *Checkout) Validate() error {
	err := Validate.Struct(co)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDCheckout, err.(validator.ValidationErrors))
	}
	return nil
}

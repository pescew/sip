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

const MsgIDRenew = "29"

var ErrInvalidRequest29 = fmt.Errorf("Invalid SIP %s request", MsgIDRenew)

// This message is used to renew an item. The ACS should respond with a Renew Response message. Either or both of the “item identifier” and “title identifier” fields must be present for the message to be useful.
type Renew struct {
	// Required:
	ThirdPartyAllowed bool
	NoBlock           bool
	TransactionDate   time.Time `validate:"required"`
	NBDueDate         time.Time `validate:"required"`
	InstitutionID     string    `validate:"required,sip"`
	PatronID          string    `validate:"required,sip"`

	// Optional:
	PatronPassword   string `validate:"sip"`
	ItemID           string `validate:"sip"`
	TitleID          string `validate:"sip"`
	TerminalPassword string `validate:"sip"`
	ItemProperties   string `validate:"sip"`
	FeeAcknowledged  bool
}

func (rn *Renew) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDRenew)

	msg.WriteString(utils.YorN(rn.ThirdPartyAllowed))
	msg.WriteString(utils.YorN(rn.NoBlock))

	msg.WriteString(rn.TransactionDate.Format(utils.SIPDateFormat))
	msg.WriteString(rn.NBDueDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", rn.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", rn.PatronID, delimiter)

	if rn.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", rn.PatronPassword, delimiter)
	}

	if rn.ItemID != "" {
		fmt.Fprintf(&msg, "AB%s%c", rn.ItemID, delimiter)
	}

	if rn.TitleID != "" {
		fmt.Fprintf(&msg, "AJ%s%c", rn.TitleID, delimiter)
	}

	if rn.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", rn.TerminalPassword, delimiter)
	}

	if rn.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", rn.ItemProperties, delimiter)
	}

	if rn.FeeAcknowledged {
		fmt.Fprintf(&msg, "BO%s%c", utils.YorN(rn.FeeAcknowledged), delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (rn *Renew) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 46 {
		return 0, ErrInvalidRequest29
	}

	if string(runes[0:2]) != MsgIDRenew {
		return 0, ErrInvalidRequest29
	}

	codes := utils.ExtractFields(string(runes[40:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AD": "", "AB": "", "AJ": "", "AC": "", "CH": "", "BO": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	rn.ThirdPartyAllowed = utils.ParseBool(runes[2])
	rn.NoBlock = utils.ParseBool(runes[3])

	rn.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[4:22]))
	if err != nil {
		return 0, err
	}

	rn.NBDueDate, err = time.Parse(utils.SIPDateFormat, string(runes[22:40]))
	if err != nil {
		return 0, err
	}

	rn.InstitutionID = codes["AO"]
	rn.PatronID = codes["AA"]
	rn.PatronPassword = codes["AD"]
	rn.ItemID = codes["AB"]
	rn.TitleID = codes["AJ"]
	rn.TerminalPassword = codes["AC"]
	rn.ItemProperties = codes["CH"]

	if utf8.RuneCountInString(codes["BO"]) > 0 {
		rn.FeeAcknowledged = utils.ParseBool([]rune(codes["BO"])[0])
	}

	rn.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (rn *Renew) Validate() error {
	err := Validate.Struct(rn)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDRenew, err.(validator.ValidationErrors))
	}

	if rn.ItemID == "" && rn.TitleID == "" {
		return fmt.Errorf("%v: one of ItemID or TitleID required.", ErrInvalidRequest29)
	}

	return nil
}

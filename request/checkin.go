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

const MsgIDCheckin = "09"

var ErrInvalidRequest09 = fmt.Errorf("Invalid SIP %s request", MsgIDCheckin)

// This message is used by the SC to request to check in an item, and also to cancel a Checkout request that did not successfully complete. The ACS must respond to this command with a Checkin Response message.
type Checkin struct {
	// Required:
	NoBlock          bool
	TransactionDate  time.Time `validate:"required"`
	ReturnDate       time.Time `validate:"required"`
	CurrentLocation  string    `validate:"required,sip"`
	InstitutionID    string    `validate:"required,sip"`
	ItemID           string    `validate:"required,sip"`
	TerminalPassword string    `validate:"sip"`

	// Optional:
	ItemProperties string `validate:"sip"`
	Cancel         bool
}

func (ci *Checkin) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDCheckin)

	msg.WriteString(utils.YorN(ci.NoBlock))

	msg.WriteString(ci.TransactionDate.Format(utils.SIPDateFormat))
	msg.WriteString(ci.ReturnDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AP%s%c", ci.CurrentLocation, delimiter)
	fmt.Fprintf(&msg, "AO%s%c", ci.InstitutionID, delimiter)

	fmt.Fprintf(&msg, "AB%s%c", ci.ItemID, delimiter)
	fmt.Fprintf(&msg, "AC%s%c", ci.TerminalPassword, delimiter)

	if ci.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", ci.ItemProperties, delimiter)
	}

	if ci.Cancel {
		fmt.Fprintf(&msg, "BI%s%c", utils.YorN(ci.Cancel), delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ci *Checkin) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 51 {
		return 0, ErrInvalidRequest09
	}

	if string(runes[0:2]) != MsgIDCheckin {
		return 0, ErrInvalidRequest09
	}

	codes := utils.ExtractFields(string(runes[39:]), delimiter, map[string]string{"AY": "", "AP": "", "AO": "", "AB": "", "AC": "", "CH": "", "BI": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	ci.NoBlock = utils.ParseBool(runes[2])

	ci.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	ci.ReturnDate, err = time.Parse(utils.SIPDateFormat, string(runes[21:39]))
	if err != nil {
		return 0, err
	}

	ci.CurrentLocation = codes["AP"]
	ci.InstitutionID = codes["AO"]
	ci.ItemID = codes["AB"]
	ci.TerminalPassword = codes["AC"]

	ci.ItemProperties = codes["CH"]

	if utf8.RuneCountInString(codes["BI"]) > 0 {
		ci.Cancel = utils.ParseBool([]rune(codes["BI"])[0])
	}

	err = ci.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (ci *Checkin) Validate() error {
	err := Validate.Struct(ci)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDCheckin, err.(validator.ValidationErrors))
	}
	return nil
}

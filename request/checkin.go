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

var ErrInvalidRequest09 = fmt.Errorf("Invalid SIP %s request", types.ReqCheckin.String())

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

	SeqNum int `validate:"min=0,max=9"`
}

func (ci *Checkin) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqCheckin.ID())

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

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", ci.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (ci *Checkin) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 51 {
		return ErrInvalidRequest09
	}

	if string(runes[0:2]) != types.ReqCheckin.ID() {
		return ErrInvalidRequest09
	}

	codes := utils.ExtractFields(string(runes[39:]), delimiter, map[string]string{"AY": "", "AP": "", "AO": "", "AB": "", "AC": "", "CH": "", "BI": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ci.SeqNum = 0
	} else {
		ci.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ci.SeqNum = 0
		}
	}

	ci.NoBlock = utils.ParseBool(runes[2])

	ci.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return err
	}

	ci.ReturnDate, err = time.Parse(utils.SIPDateFormat, string(runes[21:39]))
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (ci *Checkin) Validate() error {
	err := Validate.Struct(ci)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqCheckin.String(), err.(validator.ValidationErrors))
	}
	return nil
}

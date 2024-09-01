package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDCheckin = "10"

var ErrInvalidResponse10 = fmt.Errorf("Invalid SIP %s response", MsgIDCheckin)

// This message must be sent by the ACS in response to a SC Checkin message.
type Checkin struct {
	// Required Fields:
	Ok                bool
	Resensitize       bool
	MagneticMedia     bool
	Alert             bool
	TransactionDate   time.Time `validate:"required"`
	InstitutionID     string    `validate:"required,sip"`
	ItemID            string    `validate:"required,sip"`
	PermanentLocation string    `validate:"required,sip"`

	// Optional Fields:
	TitleID        string `validate:"sip"`
	SortBin        string `validate:"sip"`
	PatronID       string `validate:"sip"`
	MediaType      string `validate:"sip,max=3"`
	ItemProperties string `validate:"sip"`
	ScreenMessage  string `validate:"sip"`
	PrintLine      string `validate:"sip"`
}

func (ci *Checkin) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDCheckin)

	msg.WriteString(utils.ZeroOrOne(ci.Ok))
	msg.WriteString(utils.YorN(ci.Resensitize))
	msg.WriteString(utils.YorN(ci.MagneticMedia))
	msg.WriteString(utils.YorN(ci.Alert))
	msg.WriteString(ci.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", ci.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", ci.ItemID, delimiter)
	fmt.Fprintf(&msg, "AQ%s%c", ci.PermanentLocation, delimiter)

	if ci.TitleID != "" {
		fmt.Fprintf(&msg, "AJ%s%c", ci.TitleID, delimiter)
	}

	if ci.SortBin != "" {
		fmt.Fprintf(&msg, "CL%s%c", ci.SortBin, delimiter)
	}

	if ci.PatronID != "" {
		fmt.Fprintf(&msg, "AA%s%c", ci.PatronID, delimiter)
	}

	if utf8.RuneCountInString(ci.MediaType) == 3 {
		fmt.Fprintf(&msg, "CK%s%c", ci.MediaType, delimiter)
	}

	if ci.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", ci.ItemProperties, delimiter)
	}

	if ci.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", ci.ScreenMessage, delimiter)
	}

	if ci.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", ci.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ci *Checkin) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 30 {
		return 0, ErrInvalidResponse10
	}

	if string(runes[0:2]) != MsgIDCheckin {
		return 0, ErrInvalidResponse10
	}

	codes := utils.ExtractFields(string(runes[24:]), delimiter, map[string]string{"AY": "", "AO": "", "AB": "", "AQ": "", "AJ": "", "CL": "", "AA": "", "CK": "", "CH": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	ci.Ok = utils.ParseBool(runes[2])
	ci.Resensitize = utils.ParseBool(runes[3])
	ci.MagneticMedia = utils.ParseBool(runes[4])
	ci.Alert = utils.ParseBool(runes[5])

	ci.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[6:24]))
	if err != nil {
		return 0, err
	}

	ci.InstitutionID = codes["AO"]
	ci.ItemID = codes["AB"]
	ci.PermanentLocation = codes["AQ"]
	ci.TitleID = codes["AJ"]
	ci.SortBin = codes["CL"]
	ci.PatronID = codes["AA"]

	if utf8.RuneCountInString(codes["CK"]) == 3 {
		ci.MediaType = codes["CK"]
	}

	ci.ItemProperties = codes["CH"]
	ci.ScreenMessage = codes["AF"]
	ci.PrintLine = codes["AG"]

	err = ci.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (ci *Checkin) Validate() error {
	if ci.MediaType != "" && utf8.RuneCountInString(ci.MediaType) != 3 {
		return fmt.Errorf("invalid SIP %s response did not pass validation: MediaType must be 3 chars", MsgIDCheckin)
	}

	err := Validate.Struct(ci)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDCheckin, err.(validator.ValidationErrors))
	}
	return nil
}

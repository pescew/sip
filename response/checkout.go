package response

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

var ErrInvalidResponse12 = fmt.Errorf("Invalid SIP %s", types.RespCheckout.String())

// This message must be sent by the ACS in response to a Checkout message from the SC.
type Checkout struct {
	// Required Fields:
	Ok              bool
	RenewalOk       bool
	MagneticMedia   bool
	Desensitize     bool
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`
	ItemID          string    `validate:"required,sip"`
	TitleID         string    `validate:"sip"`
	DueDate         string    `validate:"required,sip"`

	// Optional Fields:
	FeeType         int `validate:"min=0,max=99"`
	SecurityInhibit bool
	CurrencyType    string `validate:"sip,max=3"`
	FeeAmount       string `validate:"sip"`
	MediaType       string `validate:"sip,max=3"`
	ItemProperties  string `validate:"sip"`
	TransactionID   string `validate:"sip"`
	ScreenMessage   string `validate:"sip"`
	PrintLine       string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (co *Checkout) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespCheckout.ID())
	msg.WriteString(utils.ZeroOrOne(co.Ok))
	msg.WriteString(utils.YorN(co.RenewalOk))
	msg.WriteString(utils.YorN(co.MagneticMedia))
	msg.WriteString(utils.YorN(co.Desensitize))
	msg.WriteString(co.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", co.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", co.PatronID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", co.ItemID, delimiter)
	fmt.Fprintf(&msg, "AJ%s%c", co.TitleID, delimiter)
	fmt.Fprintf(&msg, "AH%s%c", co.DueDate, delimiter)

	if co.FeeType > 0 {
		fmt.Fprintf(&msg, "BT%02d%c", co.FeeType, delimiter)
	}

	fmt.Fprintf(&msg, "CI%s%c", utils.YorN(co.SecurityInhibit), delimiter)

	if utf8.RuneCountInString(co.CurrencyType) == 3 {
		fmt.Fprintf(&msg, "BH%s%c", co.CurrencyType, delimiter)
	}

	if co.FeeAmount != "" {
		fmt.Fprintf(&msg, "BV%s%c", co.FeeAmount, delimiter)
	}

	if utf8.RuneCountInString(co.MediaType) == 3 {
		fmt.Fprintf(&msg, "CK%s%c", co.MediaType, delimiter)
	}

	if co.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", co.ItemProperties, delimiter)
	}

	if co.TransactionID != "" {
		fmt.Fprintf(&msg, "BK%s%c", co.TransactionID, delimiter)
	}

	if co.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", co.ScreenMessage, delimiter)
	}

	if co.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", co.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", co.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (co *Checkout) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 39 {
		return ErrInvalidResponse12
	}

	if string(runes[0:2]) != types.RespCheckout.ID() {
		return ErrInvalidResponse12
	}

	codes := utils.ExtractFields(string(runes[24:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AB": "", "AJ": "", "AH": "", "BT": "", "CI": "", "BH": "", "BV": "", "CK": "", "CH": "", "BK": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		co.SeqNum = 0
	} else {
		co.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			co.SeqNum = 0
		}
	}

	co.Ok = utils.ParseBool(runes[2])
	co.RenewalOk = utils.ParseBool(runes[3])
	co.MagneticMedia = utils.ParseBool(runes[4])
	co.Desensitize = utils.ParseBool(runes[5])

	co.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[6:24]))
	if err != nil {
		return err
	}

	co.InstitutionID = codes["AO"]
	co.PatronID = codes["AA"]
	co.ItemID = codes["AB"]
	co.TitleID = codes["AJ"]
	co.DueDate = codes["AH"]

	if codes["BT"] != "" {
		co.FeeType, err = strconv.Atoi(codes["BT"])
		if err != nil {
			return err
		}
	}

	if codes["CI"] != "" {
		co.SecurityInhibit = utils.ParseBool([]rune(codes["CI"])[0])
	}

	if utf8.RuneCountInString(codes["BH"]) == 3 {
		co.CurrencyType = codes["BH"]
	}

	co.FeeAmount = codes["BV"]

	if utf8.RuneCountInString(codes["CK"]) == 3 {
		co.MediaType = codes["CK"]
	}

	co.ItemProperties = codes["CH"]
	co.TransactionID = codes["BK"]
	co.ScreenMessage = codes["AF"]
	co.PrintLine = codes["AG"]

	err = co.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (co *Checkout) Validate() error {
	if co.CurrencyType != "" && utf8.RuneCountInString(co.CurrencyType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: CurrencyType must be 3 chars", types.RespCheckout.String())
	}

	if co.MediaType != "" && utf8.RuneCountInString(co.MediaType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: MediaType must be 3 chars", types.RespCheckout.String())
	}

	err := Validate.Struct(co)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespCheckout.String(), err.(validator.ValidationErrors))
	}
	return nil
}

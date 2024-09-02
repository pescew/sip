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

var ErrInvalidResponse30 = fmt.Errorf("Invalid SIP %s", types.RespRenew.String())

// This message must be sent by the ACS in response to a Renew message by the SC.
type Renew struct {
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

func (rn *Renew) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespRenew.ID())
	msg.WriteString(utils.ZeroOrOne(rn.Ok))
	msg.WriteString(utils.YorN(rn.RenewalOk))
	msg.WriteString(utils.YorN(rn.MagneticMedia))
	msg.WriteString(utils.YorN(rn.Desensitize))
	msg.WriteString(rn.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", rn.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", rn.PatronID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", rn.ItemID, delimiter)
	fmt.Fprintf(&msg, "AJ%s%c", rn.TitleID, delimiter)
	fmt.Fprintf(&msg, "AH%s%c", rn.DueDate, delimiter)

	if rn.FeeType > 0 {
		fmt.Fprintf(&msg, "BT%02d%c", rn.FeeType, delimiter)
	}

	fmt.Fprintf(&msg, "CI%s%c", utils.YorN(rn.SecurityInhibit), delimiter)

	if utf8.RuneCountInString(rn.CurrencyType) == 3 {
		fmt.Fprintf(&msg, "BH%s%c", rn.CurrencyType, delimiter)
	}

	if rn.FeeAmount != "" {
		fmt.Fprintf(&msg, "BV%s%c", rn.FeeAmount, delimiter)
	}

	if utf8.RuneCountInString(rn.MediaType) == 3 {
		fmt.Fprintf(&msg, "CK%s%c", rn.MediaType, delimiter)
	}

	if rn.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", rn.ItemProperties, delimiter)
	}

	if rn.TransactionID != "" {
		fmt.Fprintf(&msg, "BK%s%c", rn.TransactionID, delimiter)
	}

	if rn.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", rn.ScreenMessage, delimiter)
	}

	if rn.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", rn.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", rn.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (rn *Renew) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 39 {
		return ErrInvalidResponse30
	}

	if string(runes[0:2]) != types.RespRenew.ID() {
		return ErrInvalidResponse30
	}

	codes := utils.ExtractFields(string(runes[24:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AB": "", "AJ": "", "AH": "", "BT": "", "CI": "", "BH": "", "BV": "", "CK": "", "CH": "", "BK": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		rn.SeqNum = 0
	} else {
		rn.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			rn.SeqNum = 0
		}
	}

	rn.Ok = utils.ParseBool(runes[2])
	rn.RenewalOk = utils.ParseBool(runes[3])
	rn.MagneticMedia = utils.ParseBool(runes[4])
	rn.Desensitize = utils.ParseBool(runes[5])

	rn.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[6:24]))
	if err != nil {
		return err
	}

	rn.InstitutionID = codes["AO"]
	rn.PatronID = codes["AA"]
	rn.ItemID = codes["AB"]
	rn.TitleID = codes["AJ"]
	rn.DueDate = codes["AH"]

	if codes["BT"] != "" {
		rn.FeeType, err = strconv.Atoi(codes["BT"])
		if err != nil {
			return err
		}
	}

	if codes["CI"] != "" {
		rn.SecurityInhibit = utils.ParseBool([]rune(codes["CI"])[0])
	}

	if utf8.RuneCountInString(codes["BH"]) == 3 {
		rn.CurrencyType = codes["BH"]
	}

	rn.FeeAmount = codes["BV"]

	if utf8.RuneCountInString(codes["CK"]) == 3 {
		rn.MediaType = codes["CK"]
	}

	rn.ItemProperties = codes["CH"]
	rn.TransactionID = codes["BK"]
	rn.ScreenMessage = codes["AF"]
	rn.PrintLine = codes["AG"]

	err = rn.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (rn *Renew) Validate() error {
	if rn.CurrencyType != "" && utf8.RuneCountInString(rn.CurrencyType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: CurrencyType must be 3 chars", types.RespRenew.String())
	}

	if rn.MediaType != "" && utf8.RuneCountInString(rn.MediaType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: MediaType must be 3 chars", types.RespRenew.String())
	}

	err := Validate.Struct(rn)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespRenew.String(), err.(validator.ValidationErrors))
	}
	return nil
}

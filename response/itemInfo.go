package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse18 = fmt.Errorf("Invalid SIP %s", fields.RespItemInfo.String())

// The ACS must send this message in response to the Item Information message.
type ItemInfo struct {
	// Required Fields:
	CirculationStatus fields.CirculationStatus `validate:"min=0,max=99"`
	SecurityMarker    fields.SecurityMarker    `validate:"min=0,max=99"`
	FeeType           fields.FeeType           `validate:"min=1,max=99"`
	TransactionDate   time.Time                `validate:"required"`

	// Optional Fields:
	HoldQueueLength int    `validate:"min=-1"`
	DueDate         string `validate:"required,sip"`
	RecallDate      time.Time
	HoldPickupDate  time.Time

	// Required Fields:
	ItemID  string `validate:"required,sip"`
	TitleID string `validate:"sip"`

	// Optional Fields:
	Owner             string           `validate:"sip"`
	CurrencyType      string           `validate:"sip,max=3"`
	FeeAmount         string           `validate:"sip"`
	MediaType         fields.MediaType `validate:"sip,max=3"`
	PermanentLocation string           `validate:"sip"`
	CurrentLocation   string           `validate:"sip"`
	ItemProperties    string           `validate:"sip"`
	ScreenMessage     string           `validate:"sip"`
	PrintLine         string           `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (ii *ItemInfo) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(fields.RespItemInfo.ID())

	msg.WriteString(ii.CirculationStatus.ID())
	msg.WriteString(ii.SecurityMarker.ID())
	msg.WriteString(ii.FeeType.ID())
	msg.WriteString(ii.TransactionDate.Format(utils.SIPDateFormat))

	if ii.HoldQueueLength != -1 {
		fmt.Fprintf(&msg, "CF%d%c", ii.HoldQueueLength, delimiter)
	}

	if ii.DueDate != "" {
		fmt.Fprintf(&msg, "AH%s%c", ii.DueDate, delimiter)
	}

	if !ii.RecallDate.IsZero() {
		fmt.Fprintf(&msg, "CJ%s%c", ii.RecallDate.Format(utils.SIPDateFormat), delimiter)
	}

	if !ii.HoldPickupDate.IsZero() {
		fmt.Fprintf(&msg, "CM%s%c", ii.HoldPickupDate.Format(utils.SIPDateFormat), delimiter)
	}

	fmt.Fprintf(&msg, "AB%s%c", ii.ItemID, delimiter)
	fmt.Fprintf(&msg, "AJ%s%c", ii.TitleID, delimiter)

	if ii.Owner != "" {
		fmt.Fprintf(&msg, "BG%s%c", ii.Owner, delimiter)
	}

	if utf8.RuneCountInString(ii.CurrencyType) == 3 {
		fmt.Fprintf(&msg, "BH%s%c", ii.CurrencyType, delimiter)
	}

	if ii.FeeAmount != "" {
		fmt.Fprintf(&msg, "BV%s%c", ii.FeeAmount, delimiter)
	}

	if utf8.RuneCountInString(ii.MediaType.ID()) == 3 {
		fmt.Fprintf(&msg, "CK%s%c", ii.MediaType.ID(), delimiter)
	}

	if ii.PermanentLocation != "" {
		fmt.Fprintf(&msg, "AQ%s%c", ii.PermanentLocation, delimiter)
	}

	if ii.CurrentLocation != "" {
		fmt.Fprintf(&msg, "AP%s%c", ii.CurrentLocation, delimiter)
	}

	if ii.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", ii.ItemProperties, delimiter)
	}

	if ii.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", ii.ScreenMessage, delimiter)
	}

	if ii.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", ii.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", ii.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (ii *ItemInfo) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 32 {
		return ErrInvalidResponse18
	}

	if string(runes[0:2]) != fields.RespItemInfo.ID() {
		return ErrInvalidResponse18
	}

	codes := utils.ExtractFields(string(runes[26:]), delimiter, map[string]string{"AY": "", "CF": "", "AH": "", "CJ": "", "CM": "", "AB": "", "AJ": "", "BG": "", "BH": "", "BV": "", "CK": "", "AQ": "", "AP": "", "CH": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ii.SeqNum = 0
	} else {
		ii.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ii.SeqNum = 0
		}
	}

	circStatus, err := strconv.Atoi(string(runes[2:4]))
	if err != nil {
		return err
	}
	ii.CirculationStatus = fields.CirculationStatus(circStatus)

	securityMarker, err := strconv.Atoi(string(runes[4:6]))
	if err != nil {
		return err
	}
	ii.SecurityMarker = fields.SecurityMarker(securityMarker)

	feeType, err := strconv.Atoi(string(runes[6:8]))
	if err != nil {
		return err
	}
	ii.FeeType = fields.FeeType(feeType)

	ii.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[8:26]))
	if err != nil {
		return err
	}

	if codes["CF"] != "" {
		ii.HoldQueueLength, err = strconv.Atoi(codes["CF"])
		if err != nil {
			return err
		}
	} else {
		ii.HoldQueueLength = -1
	}

	ii.DueDate = codes["AH"]

	if codes["CJ"] != "" {
		ii.RecallDate, err = time.Parse(utils.SIPDateFormat, codes["CJ"])
		if err != nil {
			return err
		}
	}

	if codes["CM"] != "" {
		ii.HoldPickupDate, err = time.Parse(utils.SIPDateFormat, codes["CM"])
		if err != nil {
			return err
		}
	}

	ii.ItemID = codes["AB"]
	ii.TitleID = codes["AJ"]
	ii.Owner = codes["BG"]

	if utf8.RuneCountInString(codes["BH"]) == 3 {
		ii.CurrencyType = codes["BH"]
	}

	ii.FeeAmount = codes["BV"]

	if utf8.RuneCountInString(codes["CK"]) == 3 {
		ii.MediaType = fields.MediaType(codes["CK"])
	}

	ii.PermanentLocation = codes["AQ"]
	ii.CurrentLocation = codes["AP"]
	ii.ItemProperties = codes["CH"]
	ii.ScreenMessage = codes["AF"]
	ii.PrintLine = codes["AG"]

	err = ii.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (ii *ItemInfo) Validate() error {
	if ii.CurrencyType != "" && utf8.RuneCountInString(ii.CurrencyType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: CurrencyType must be 3 chars", fields.RespItemInfo.String())
	}

	if ii.MediaType.ID() != "" && utf8.RuneCountInString(ii.MediaType.ID()) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: MediaType must be 3 chars", fields.RespItemInfo.String())
	}

	err := Validate.Struct(ii)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", fields.RespItemInfo.String(), err.(validator.ValidationErrors))
	}
	return nil
}

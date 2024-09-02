package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse64 = fmt.Errorf("Invalid SIP %s", types.RespPatronInfo.String())

// The ACS must send this message in response to the Patron Information message.
type PatronInfo struct {
	// Required Fields:
	PatronStatus          fields.PatronStatus `validate:"required"`
	Language              int                 `validate:"min=0,max=999"`
	TransactionDate       time.Time           `validate:"required"`
	HoldItemsCount        int                 `validate:"min=0,max=9999"`
	OverdueItemsCount     int                 `validate:"min=0,max=9999"`
	ChargedItemsCount     int                 `validate:"min=0,max=9999"`
	FineItemsCount        int                 `validate:"min=0,max=9999"`
	RecallItemsCount      int                 `validate:"min=0,max=9999"`
	UnavailableHoldsCount int                 `validate:"min=0,max=9999"`
	InstitutionID         string              `validate:"required,sip"`
	PatronID              string              `validate:"required,sip"`
	PatronName            string              `validate:"sip"`

	// Optional Fields:
	HoldItemsLimit      int `validate:"min=0,max=9999"`
	OverdueItemsLimit   int `validate:"min=0,max=9999"`
	ChargedItemsLimit   int `validate:"min=0,max=9999"`
	ValidPatron         bool
	ValidPatronPassword bool
	CurrencyType        string   `validate:"sip,max=3"`
	FeeAmount           string   `validate:"sip"`
	FeeLimit            string   `validate:"sip"`
	HoldItems           []string `validate:"dive,sip"`
	OverdueItems        []string `validate:"dive,sip"`
	ChargedItems        []string `validate:"dive,sip"`
	FineItems           []string `validate:"dive,sip"`
	RecallItems         []string `validate:"dive,sip"`
	UnavailHoldItems    []string `validate:"dive,sip"`
	HomeAddress         string   `validate:"sip"`
	EmailAddress        string   `validate:"sip"`
	HomePhone           string   `validate:"sip"`
	ScreenMessage       string   `validate:"sip"`
	PrintLine           string   `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (pi *PatronInfo) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespPatronInfo.ID())
	msg.WriteString(pi.PatronStatus.Marshal())

	fmt.Fprintf(&msg, "%03d", pi.Language)
	msg.WriteString(pi.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "%04d", pi.HoldItemsCount)
	fmt.Fprintf(&msg, "%04d", pi.OverdueItemsCount)
	fmt.Fprintf(&msg, "%04d", pi.ChargedItemsCount)
	fmt.Fprintf(&msg, "%04d", pi.FineItemsCount)
	fmt.Fprintf(&msg, "%04d", pi.RecallItemsCount)
	fmt.Fprintf(&msg, "%04d", pi.UnavailableHoldsCount)

	fmt.Fprintf(&msg, "AO%s%c", pi.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", pi.PatronID, delimiter)
	fmt.Fprintf(&msg, "AE%s%c", pi.PatronName, delimiter)

	if pi.HoldItemsLimit > 0 {
		fmt.Fprintf(&msg, "BZ%04d%c", pi.HoldItemsLimit, delimiter)
	}

	if pi.OverdueItemsLimit > 0 {
		fmt.Fprintf(&msg, "CA%04d%c", pi.OverdueItemsLimit, delimiter)
	}

	if pi.ChargedItemsLimit > 0 {
		fmt.Fprintf(&msg, "CB%04d%c", pi.ChargedItemsLimit, delimiter)
	}

	//todo: check if patron is valid?
	fmt.Fprintf(&msg, "BL%s%c", utils.YorN(pi.ValidPatron), delimiter)

	fmt.Fprintf(&msg, "CQ%s%c", utils.YorN(pi.ValidPatronPassword), delimiter)

	if utf8.RuneCountInString(pi.CurrencyType) == 3 {
		fmt.Fprintf(&msg, "BH%s%c", pi.CurrencyType, delimiter)
	}

	if pi.FeeAmount != "" {
		fmt.Fprintf(&msg, "BV%s%c", pi.FeeAmount, delimiter)
	}

	if pi.FeeLimit != "" {
		fmt.Fprintf(&msg, "CC%s%c", pi.FeeLimit, delimiter)
	}

	for _, holdItem := range pi.HoldItems {
		if holdItem != "" {
			fmt.Fprintf(&msg, "AS%s%c", holdItem, delimiter)
		}
	}
	for _, overdueItem := range pi.OverdueItems {
		if overdueItem != "" {
			fmt.Fprintf(&msg, "AT%s%c", overdueItem, delimiter)
		}
	}
	for _, chargedItem := range pi.ChargedItems {
		if chargedItem != "" {
			fmt.Fprintf(&msg, "AU%s%c", chargedItem, delimiter)
		}
	}
	for _, fineItem := range pi.FineItems {
		if fineItem != "" {
			fmt.Fprintf(&msg, "AV%s%c", fineItem, delimiter)
		}
	}
	for _, recallItem := range pi.RecallItems {
		if recallItem != "" {
			fmt.Fprintf(&msg, "BU%s%c", recallItem, delimiter)
		}
	}
	for _, unavailHoldItem := range pi.UnavailHoldItems {
		if unavailHoldItem != "" {
			fmt.Fprintf(&msg, "CD%s%c", unavailHoldItem, delimiter)
		}
	}

	if pi.HomeAddress != "" {
		fmt.Fprintf(&msg, "BD%s%c", pi.HomeAddress, delimiter)
	}

	if pi.EmailAddress != "" {
		fmt.Fprintf(&msg, "BE%s%c", pi.EmailAddress, delimiter)
	}

	if pi.HomePhone != "" {
		fmt.Fprintf(&msg, "BF%s%c", pi.HomePhone, delimiter)
	}

	if pi.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", pi.ScreenMessage, delimiter)
	}

	if pi.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", pi.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", pi.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (pi *PatronInfo) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 73 {
		return ErrInvalidResponse64
	}

	if string(runes[0:2]) != types.RespPatronInfo.ID() {
		return ErrInvalidResponse64
	}

	codes := utils.ExtractFields(string(runes[61:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AE": "", "BZ": "", "CA": "", "CB": "", "BL": "", "CQ": "", "BH": "", "BV": "", "CC": "", "BD": "", "BE": "", "BF": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		pi.SeqNum = 0
	} else {
		pi.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			pi.SeqNum = 0
		}
	}

	multiCodes := utils.ExtractMultiFields(string(runes[70:]), delimiter, map[string][]string{"AS": []string{}, "AT": []string{}, "AU": []string{}, "AV": []string{}, "BU": []string{}, "CD": []string{}})

	err = pi.PatronStatus.Unmarshal(string(runes[2:16]))
	if err != nil {
		return err
	}

	pi.Language, err = strconv.Atoi(string(runes[16:19]))
	if err != nil {
		return err
	}

	pi.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[19:37]))
	if err != nil {
		return err
	}

	pi.HoldItemsCount, err = strconv.Atoi(string(runes[37:41]))
	if err != nil {
		return err
	}

	pi.OverdueItemsCount, err = strconv.Atoi(string(runes[41:45]))
	if err != nil {
		return err
	}

	pi.ChargedItemsCount, err = strconv.Atoi(string(runes[45:49]))
	if err != nil {
		return err
	}

	pi.FineItemsCount, err = strconv.Atoi(string(runes[49:53]))
	if err != nil {
		return err
	}

	pi.RecallItemsCount, err = strconv.Atoi(string(runes[53:57]))
	if err != nil {
		return err
	}

	pi.UnavailableHoldsCount, err = strconv.Atoi(string(runes[57:61]))
	if err != nil {
		return err
	}

	pi.InstitutionID = codes["AO"]
	pi.PatronID = codes["AA"]
	pi.PatronName = codes["AE"]

	if codes["BZ"] != "" {
		pi.HoldItemsLimit, err = strconv.Atoi(codes["BZ"])
		if err != nil {
			return err
		}
	}

	if codes["CA"] != "" {
		pi.OverdueItemsLimit, err = strconv.Atoi(codes["CA"])
		if err != nil {
			return err
		}
	}

	if codes["CB"] != "" {
		pi.ChargedItemsLimit, err = strconv.Atoi(codes["CB"])
		if err != nil {
			return err
		}
	}

	if codes["BL"] != "" {
		pi.ValidPatron = utils.ParseBool([]rune(codes["BL"])[0])
	}

	if codes["CQ"] != "" {
		pi.ValidPatronPassword = utils.ParseBool([]rune(codes["CQ"])[0])
	}

	if utf8.RuneCountInString(codes["BH"]) == 3 {
		pi.CurrencyType = codes["BH"]
	}

	pi.FeeAmount = codes["BV"]
	pi.FeeLimit = codes["CC"]

	pi.HoldItems = multiCodes["AS"]
	pi.OverdueItems = multiCodes["AT"]
	pi.ChargedItems = multiCodes["AU"]
	pi.FineItems = multiCodes["AV"]
	pi.RecallItems = multiCodes["BU"]
	pi.UnavailHoldItems = multiCodes["CD"]

	pi.HomeAddress = codes["BD"]
	pi.EmailAddress = codes["BE"]
	pi.HomePhone = codes["BF"]
	pi.ScreenMessage = codes["AF"]
	pi.PrintLine = codes["AG"]

	err = pi.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (pi *PatronInfo) Validate() error {
	if pi.CurrencyType != "" && utf8.RuneCountInString(pi.CurrencyType) != 3 {
		return fmt.Errorf("invalid SIP %s did not pass validation: CurrencyType must be 3 chars", types.RespPatronInfo.String())
	}

	err := Validate.Struct(pi)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespPatronInfo.String(), err.(validator.ValidationErrors))
	}
	return nil
}

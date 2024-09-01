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

const MsgIDPatronStatus = "24"

var ErrInvalidResponse24 = fmt.Errorf("Invalid SIP %s response", MsgIDPatronStatus)

// The ACS must send this message in response to a Patron Status Request message as well as in response to a Block Patron message.
type PatronStatus struct {
	// Required Fields:
	PatronStatus    fields.PatronStatus `validate:"required"`
	Language        int                 `validate:"min=0,max=999"`
	TransactionDate time.Time           `validate:"required"`
	InstitutionID   string              `validate:"sip"`
	PatronID        string              `validate:"sip"`
	PatronName      string              `validate:"sip"`

	// Optional Fields:
	ValidPatron         bool
	ValidPatronPassword bool
	CurrencyType        string `validate:"sip,max=3"`
	FeeAmount           string `validate:"sip"`
	ScreenMessage       string `validate:"sip"`
	PrintLine           string `validate:"sip"`
}

func (ps *PatronStatus) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDPatronStatus)

	msg.WriteString(ps.PatronStatus.Marshal())
	fmt.Fprintf(&msg, "%03d", ps.Language)
	msg.WriteString(ps.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", ps.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", ps.PatronID, delimiter)
	fmt.Fprintf(&msg, "AE%s%c", ps.PatronName, delimiter)

	fmt.Fprintf(&msg, "BL%s%c", utils.YorN(ps.ValidPatron), delimiter)
	fmt.Fprintf(&msg, "CQ%s%c", utils.YorN(ps.ValidPatronPassword), delimiter)

	if utf8.RuneCountInString(ps.CurrencyType) == 3 {
		fmt.Fprintf(&msg, "BH%s%c", ps.CurrencyType, delimiter)
	}

	if ps.FeeAmount != "" {
		fmt.Fprintf(&msg, "BV%s%c", ps.FeeAmount, delimiter)
	}

	if ps.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", ps.ScreenMessage, delimiter)
	}

	if ps.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", ps.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ps *PatronStatus) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 46 {
		return 0, ErrInvalidResponse24
	}

	if string(runes[0:2]) != MsgIDPatronStatus {
		return 0, ErrInvalidResponse24
	}

	codes := utils.ExtractFields(string(runes[37:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AE": "", "BL": "", "CQ": "", "BH": "", "BV": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	err = ps.PatronStatus.Unmarshal(string(runes[2:16]))
	if err != nil {
		return 0, err
	}

	ps.Language, err = strconv.Atoi(string(runes[16:19]))
	if err != nil {
		return 0, err
	}

	ps.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[19:37]))
	if err != nil {
		return 0, err
	}

	ps.InstitutionID = codes["AO"]
	ps.PatronID = codes["AA"]
	ps.PatronName = codes["AE"]

	if codes["BL"] != "" {
		ps.ValidPatron = utils.ParseBool([]rune(codes["BL"])[0])
	}

	if codes["CQ"] != "" {
		ps.ValidPatronPassword = utils.ParseBool([]rune(codes["CQ"])[0])
	}

	if utf8.RuneCountInString(codes["BH"]) == 3 {
		ps.CurrencyType = codes["BH"]
	}

	ps.FeeAmount = codes["BV"]
	ps.ScreenMessage = codes["AF"]
	ps.PrintLine = codes["AG"]

	err = ps.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (ps *PatronStatus) Validate() error {
	if ps.CurrencyType != "" && utf8.RuneCountInString(ps.CurrencyType) != 3 {
		return fmt.Errorf("invalid SIP %s response did not pass validation: CurrencyType must be 3 chars", MsgIDRenew)
	}

	err := Validate.Struct(ps)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDPatronStatus, err.(validator.ValidationErrors))
	}
	return nil
}

package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDFeePaid = "37"

var ErrInvalidRequest37 = fmt.Errorf("Invalid SIP %s request", MsgIDFeePaid)

// This message can be used to notify the ACS that a fee has been collected from the patron. The ACS should record this information in their database and respond with a Fee Paid Response message.
type FeePaid struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	FeeType         int       `validate:"required,min=0,max=99"`
	PaymentType     int       `validate:"required,min=0,max=99"`
	CurrencyType    string    `validate:"required,sip,len=3"`
	FeeAmount       string    `validate:"required,sip"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`
	FeeID            string `validate:"sip"`
	TransactionID    string `validate:"sip"`
}

func (fp *FeePaid) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDFeePaid)

	msg.WriteString(fp.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "%02d", fp.FeeType)
	fmt.Fprintf(&msg, "%02d", fp.PaymentType)
	msg.WriteString(fp.CurrencyType)
	fmt.Fprintf(&msg, "BV%s%c", fp.FeeAmount, delimiter)
	fmt.Fprintf(&msg, "AO%s%c", fp.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", fp.PatronID, delimiter)

	if fp.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", fp.TerminalPassword, delimiter)
	}

	if fp.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", fp.PatronPassword, delimiter)
	}

	if fp.FeeID != "" {
		fmt.Fprintf(&msg, "CG%s%c", fp.FeeID, delimiter)
	}

	if fp.TransactionID != "" {
		fmt.Fprintf(&msg, "BK%s%c", fp.TransactionID, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (fp *FeePaid) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 36 {
		return 0, ErrInvalidRequest37
	}

	if string(runes[0:2]) != MsgIDFeePaid {
		return 0, ErrInvalidRequest37
	}

	codes := utils.ExtractFields(string(runes[27:]), delimiter, map[string]string{"AY": "", "BV": "", "AO": "", "AA": "", "AC": "", "AD": "", "CG": "", "BK": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	fp.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
	}

	fp.FeeType, err = strconv.Atoi(string(runes[20:22]))
	if err != nil {
		return 0, err
	}

	fp.PaymentType, err = strconv.Atoi(string(runes[22:24]))
	if err != nil {
		return 0, err
	}

	fp.CurrencyType = string(runes[24:27])

	fp.FeeAmount = codes["BV"]
	fp.InstitutionID = codes["AO"]
	fp.PatronID = codes["AA"]

	fp.TerminalPassword = codes["AC"]
	fp.PatronPassword = codes["AD"]
	fp.FeeID = codes["CG"]
	fp.TransactionID = codes["BK"]

	fp.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (fp *FeePaid) Validate() error {
	err := Validate.Struct(fp)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDFeePaid, err.(validator.ValidationErrors))
	}
	return nil
}

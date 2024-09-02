package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest37 = fmt.Errorf("Invalid SIP %s request", types.ReqFeePaid.String())

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

	SeqNum int `validate:"min=0,max=9"`
}

func (fp *FeePaid) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqFeePaid.ID())

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

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", fp.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (fp *FeePaid) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 36 {
		return ErrInvalidRequest37
	}

	if string(runes[0:2]) != types.ReqFeePaid.ID() {
		return ErrInvalidRequest37
	}

	codes := utils.ExtractFields(string(runes[27:]), delimiter, map[string]string{"AY": "", "BV": "", "AO": "", "AA": "", "AC": "", "AD": "", "CG": "", "BK": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		fp.SeqNum = 0
	} else {
		fp.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			fp.SeqNum = 0
		}
	}

	fp.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	fp.FeeType, err = strconv.Atoi(string(runes[20:22]))
	if err != nil {
		return err
	}

	fp.PaymentType, err = strconv.Atoi(string(runes[22:24]))
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (fp *FeePaid) Validate() error {
	err := Validate.Struct(fp)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqFeePaid.String(), err.(validator.ValidationErrors))
	}
	return nil
}

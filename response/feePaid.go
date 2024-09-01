package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDFeePaid = "38"

var ErrInvalidResponse38 = fmt.Errorf("Invalid SIP %s response", MsgIDFeePaid)

// The ACS must send this message in response to the Fee Paid message.
type FeePaid struct {
	// Required Fields:
	PaymentAccepted bool
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	PatronID        string    `validate:"required,sip"`

	// Optional Fields:
	TransactionID string `validate:"sip"`
	ScreenMessage string `validate:"sip"`
	PrintLine     string `validate:"sip"`
}

func (fp *FeePaid) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDFeePaid)
	msg.WriteString(utils.YorN(fp.PaymentAccepted))
	msg.WriteString(fp.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AO%s%c", fp.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", fp.PatronID, delimiter)

	if fp.TransactionID != "" {
		fmt.Fprintf(&msg, "BK%s%c", fp.TransactionID, delimiter)
	}

	if fp.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", fp.ScreenMessage, delimiter)
	}

	if fp.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", fp.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (fp *FeePaid) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 27 {
		return 0, ErrInvalidResponse38
	}

	if string(runes[0:2]) != MsgIDFeePaid {
		return 0, ErrInvalidResponse38
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "BK": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	fp.PaymentAccepted = utils.ParseBool(runes[2])

	fp.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	fp.InstitutionID = codes["AO"]
	fp.PatronID = codes["AA"]

	fp.TransactionID = codes["BK"]
	fp.ScreenMessage = codes["AF"]
	fp.PrintLine = codes["AG"]

	err = fp.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (fp *FeePaid) Validate() error {
	err := Validate.Struct(fp)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDFeePaid, err.(validator.ValidationErrors))
	}
	return nil
}

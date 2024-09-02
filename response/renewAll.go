package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse66 = fmt.Errorf("Invalid SIP %s", types.RespRenewAll.String())

// The ACS should send this message in response to a Renew All message from the SC.
type RenewAll struct {
	// Required:
	Ok              bool
	RenewedCount    int       `validate:"min=0,max=9999"`
	UnrenewedCount  int       `validate:"min=0,max=9999"`
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`

	// Optional:
	RenewedItems   []string `validate:"dive,sip"`
	UnrenewedItems []string `validate:"dive,sip"`
	ScreenMessage  string   `validate:"sip"`
	PrintLine      string   `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (ra *RenewAll) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespRenewAll.ID())

	msg.WriteString(utils.ZeroOrOne(ra.Ok))

	fmt.Fprintf(&msg, "%04d", ra.RenewedCount)
	fmt.Fprintf(&msg, "%04d", ra.UnrenewedCount)

	msg.WriteString(ra.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ra.InstitutionID, delimiter)

	for _, renewedItem := range ra.RenewedItems {
		if renewedItem != "" {
			fmt.Fprintf(&msg, "BM%s%c", renewedItem, delimiter)
		}
	}

	for _, unrenewedItem := range ra.UnrenewedItems {
		if unrenewedItem != "" {
			fmt.Fprintf(&msg, "BN%s%c", unrenewedItem, delimiter)
		}
	}

	if ra.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", ra.ScreenMessage, delimiter)
	}

	if ra.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", ra.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", ra.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (ra *RenewAll) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 32 {
		return ErrInvalidResponse66
	}

	if string(runes[0:2]) != types.RespRenewAll.ID() {
		return ErrInvalidResponse66
	}

	codes := utils.ExtractFields(string(runes[29:]), delimiter, map[string]string{"AY": "", "AO": "", "BM": "", "BN": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ra.SeqNum = 0
	} else {
		ra.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ra.SeqNum = 0
		}
	}

	multiCodes := utils.ExtractMultiFields(string(runes[32:]), delimiter, map[string][]string{"BM": []string{}, "BN": []string{}})

	ra.Ok = utils.ParseBool(runes[2])

	ra.RenewedCount, err = strconv.Atoi(string(runes[3:7]))
	if err != nil {
		return err
	}

	ra.UnrenewedCount, err = strconv.Atoi(string(runes[7:11]))
	if err != nil {
		return err
	}

	ra.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[11:29]))
	if err != nil {
		return err
	}

	ra.InstitutionID = codes["AO"]

	ra.RenewedItems = multiCodes["BM"]
	ra.UnrenewedItems = multiCodes["BN"]

	if codes["AF"] != "" {
		ra.ScreenMessage = codes["AF"]
	}

	if codes["AG"] != "" {
		ra.PrintLine = codes["AG"]
	}

	err = ra.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (ra *RenewAll) Validate() error {
	err := Validate.Struct(ra)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespRenewAll.String(), err.(validator.ValidationErrors))
	}
	return nil
}

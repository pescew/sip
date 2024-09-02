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

var ErrInvalidResponse20 = fmt.Errorf("Invalid SIP %s", types.RespItemStatusUpdate.String())

// The ACS must send this message in response to the Item Status Update message.
type ItemStatusUpdate struct {
	// Required Fields:
	ItemPropertiesOk bool
	TransactionDate  time.Time `validate:"required"`
	ItemID           string    `validate:"required,sip"`

	// Optional Fields:
	TitleID        string `validate:"sip"`
	ItemProperties string `validate:"sip"`
	ScreenMessage  string `validate:"sip"`
	PrintLine      string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (isu *ItemStatusUpdate) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespItemStatusUpdate.ID())

	msg.WriteString(utils.ZeroOrOne(isu.ItemPropertiesOk))
	msg.WriteString(isu.TransactionDate.Format(utils.SIPDateFormat))
	fmt.Fprintf(&msg, "AB%s%c", isu.ItemID, delimiter)

	if isu.TitleID != "" {
		fmt.Fprintf(&msg, "AJ%s%c", isu.TitleID, delimiter)
	}

	if isu.ItemProperties != "" {
		fmt.Fprintf(&msg, "CH%s%c", isu.ItemProperties, delimiter)
	}

	if isu.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", isu.ScreenMessage, delimiter)
	}

	if isu.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", isu.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", isu.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (isu *ItemStatusUpdate) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 24 {
		return ErrInvalidResponse20
	}

	if string(runes[0:2]) != types.RespItemStatusUpdate.ID() {
		return ErrInvalidResponse20
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AB": "", "AJ": "", "CH": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		isu.SeqNum = 0
	} else {
		isu.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			isu.SeqNum = 0
		}
	}

	isu.ItemPropertiesOk = utils.ParseBool(runes[2])

	isu.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return err
	}

	isu.ItemID = codes["AB"]

	isu.TitleID = codes["AJ"]
	isu.ItemProperties = codes["CH"]
	isu.ScreenMessage = codes["AF"]
	isu.PrintLine = codes["AG"]

	err = isu.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (isu *ItemStatusUpdate) Validate() error {
	err := Validate.Struct(isu)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespItemStatusUpdate.String(), err.(validator.ValidationErrors))
	}
	return nil
}

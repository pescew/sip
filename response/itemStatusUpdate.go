package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDItemStatusUpdate = "20"

var ErrInvalidResponse20 = fmt.Errorf("Invalid SIP %s response", MsgIDItemStatusUpdate)

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
}

func (isu *ItemStatusUpdate) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDItemStatusUpdate)

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

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (isu *ItemStatusUpdate) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 24 {
		return 0, ErrInvalidResponse20
	}

	if string(runes[0:2]) != MsgIDItemStatusUpdate {
		return 0, ErrInvalidResponse20
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AB": "", "AJ": "", "CH": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	isu.ItemPropertiesOk = utils.ParseBool(runes[2])

	isu.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	isu.ItemID = codes["AB"]

	isu.TitleID = codes["AJ"]
	isu.ItemProperties = codes["CH"]
	isu.ScreenMessage = codes["AF"]
	isu.PrintLine = codes["AG"]

	err = isu.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (isu *ItemStatusUpdate) Validate() error {
	err := Validate.Struct(isu)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDItemStatusUpdate, err.(validator.ValidationErrors))
	}
	return nil
}

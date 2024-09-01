package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDItemInfo = "17"

var ErrInvalidRequest17 = fmt.Errorf("Invalid SIP %s request", MsgIDItemInfo)

// This message may be used to request item information. The ACS should respond with the Item Information Response message.
type ItemInfo struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	ItemID          string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
}

func (ii *ItemInfo) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDItemInfo)

	msg.WriteString(ii.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ii.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", ii.ItemID, delimiter)

	if ii.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", ii.TerminalPassword, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (ii *ItemInfo) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 29 {
		return 0, ErrInvalidRequest17
	}

	if string(runes[0:2]) != MsgIDItemInfo {
		return 0, ErrInvalidRequest17
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AB": "", "AC": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	ii.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
	}

	ii.InstitutionID = codes["AO"]
	ii.ItemID = codes["AB"]
	ii.TerminalPassword = codes["AC"]

	ii.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (ii *ItemInfo) Validate() error {
	err := Validate.Struct(ii)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDItemInfo, err.(validator.ValidationErrors))
	}
	return nil
}

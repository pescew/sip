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

var ErrInvalidRequest17 = fmt.Errorf("Invalid SIP %s request", types.ReqItemInfo.String())

// This message may be used to request item information. The ACS should respond with the Item Information Response message.
type ItemInfo struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	ItemID          string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (ii *ItemInfo) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqItemInfo.ID())

	msg.WriteString(ii.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", ii.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", ii.ItemID, delimiter)

	if ii.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", ii.TerminalPassword, delimiter)
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

	if len(runes) < 29 {
		return ErrInvalidRequest17
	}

	if string(runes[0:2]) != types.ReqItemInfo.ID() {
		return ErrInvalidRequest17
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AB": "", "AC": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		ii.SeqNum = 0
	} else {
		ii.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			ii.SeqNum = 0
		}
	}

	ii.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	ii.InstitutionID = codes["AO"]
	ii.ItemID = codes["AB"]
	ii.TerminalPassword = codes["AC"]

	ii.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (ii *ItemInfo) Validate() error {
	err := Validate.Struct(ii)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqItemInfo.String(), err.(validator.ValidationErrors))
	}
	return nil
}

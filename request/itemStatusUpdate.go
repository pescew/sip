package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDItemStatusUpdate = "19"

var ErrInvalidRequest19 = fmt.Errorf("Invalid SIP %s request", MsgIDItemStatusUpdate)

// This message can be used to send item information to the ACS, without having to do a Checkout or Checkin operation. The item properties could be stored on the ACS’s database. The ACS should respond with an Item Status Update Response message.
type ItemStatusUpdate struct {
	// Required:
	TransactionDate time.Time `validate:"required"`
	InstitutionID   string    `validate:"required,sip"`
	ItemID          string    `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`

	// Required:
	ItemProperties string `validate:"required,sip"`
}

func (isu *ItemStatusUpdate) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDItemStatusUpdate)

	msg.WriteString(isu.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", isu.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", isu.ItemID, delimiter)

	if isu.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", isu.TerminalPassword, delimiter)
	}

	fmt.Fprintf(&msg, "CH%s%c", isu.ItemProperties, delimiter)

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (isu *ItemStatusUpdate) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 29 {
		return 0, ErrInvalidRequest19
	}

	if string(runes[0:2]) != MsgIDItemStatusUpdate {
		return 0, ErrInvalidRequest19
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AB": "", "AC": "", "CH": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	isu.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return 0, err
	}

	isu.InstitutionID = codes["AO"]
	isu.ItemID = codes["AB"]
	isu.TerminalPassword = codes["AC"]
	isu.ItemProperties = codes["CH"]

	isu.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (isu *ItemStatusUpdate) Validate() error {
	err := Validate.Struct(isu)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDItemStatusUpdate, err.(validator.ValidationErrors))
	}
	return nil
}

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

var ErrInvalidRequest19 = fmt.Errorf("Invalid SIP %s request", types.ReqItemStatusUpdate.String())

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

	SeqNum int `validate:"min=0,max=9"`
}

func (isu *ItemStatusUpdate) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqItemStatusUpdate.ID())

	msg.WriteString(isu.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", isu.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AB%s%c", isu.ItemID, delimiter)

	if isu.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", isu.TerminalPassword, delimiter)
	}

	fmt.Fprintf(&msg, "CH%s%c", isu.ItemProperties, delimiter)

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

	if len(runes) < 29 {
		return ErrInvalidRequest19
	}

	if string(runes[0:2]) != types.ReqItemStatusUpdate.ID() {
		return ErrInvalidRequest19
	}

	codes := utils.ExtractFields(string(runes[20:]), delimiter, map[string]string{"AY": "", "AO": "", "AB": "", "AC": "", "CH": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		isu.SeqNum = 0
	} else {
		isu.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			isu.SeqNum = 0
		}
	}

	isu.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[2:20]))
	if err != nil {
		return err
	}

	isu.InstitutionID = codes["AO"]
	isu.ItemID = codes["AB"]
	isu.TerminalPassword = codes["AC"]
	isu.ItemProperties = codes["CH"]

	isu.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (isu *ItemStatusUpdate) Validate() error {
	err := Validate.Struct(isu)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqItemStatusUpdate.String(), err.(validator.ValidationErrors))
	}
	return nil
}

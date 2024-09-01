package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDHold = "15"

var ErrInvalidRequest15 = fmt.Errorf("Invalid SIP %s request", MsgIDHold)

// This message is used to create, modify, or delete a hold. The ACS should respond with a Hold Response message. Either or both of the “item identifier” and “title identifier” fields must be present for the message to be useful.
type Hold struct {
	// Required:
	HoldMode        string    `validate:"required,len=1,oneof=+ - *"`
	TransactionDate time.Time `validate:"required"`

	// Optional:
	ExpirationDate time.Time
	PickupLocation string `validate:"sip"`
	HoldType       int    `validate:"min=0,max=9"`

	// Required:
	InstitutionID string `validate:"required,sip"`
	PatronID      string `validate:"required,sip"`

	// Optional:
	PatronPassword   string `validate:"sip"`
	ItemID           string `validate:"sip"`
	TitleID          string `validate:"sip"`
	TerminalPassword string `validate:"sip"`
	FeeAcknowledged  bool
}

func (h *Hold) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDHold)

	msg.WriteString(h.HoldMode)

	msg.WriteString(h.TransactionDate.Format(utils.SIPDateFormat))

	if !h.ExpirationDate.IsZero() {
		fmt.Fprintf(&msg, "BW%s%c", h.ExpirationDate.Format(utils.SIPDateFormat), delimiter)
	}

	if h.PickupLocation != "" {
		fmt.Fprintf(&msg, "BS%s%c", h.PickupLocation, delimiter)
	}

	if h.HoldType > 0 {
		fmt.Fprintf(&msg, "BY%s%c", strconv.Itoa(h.HoldType), delimiter)
	}

	fmt.Fprintf(&msg, "AO%s%c", h.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", h.PatronID, delimiter)

	if h.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", h.PatronPassword, delimiter)
	}

	if h.ItemID != "" {
		fmt.Fprintf(&msg, "AB%s%c", h.ItemID, delimiter)
	}

	if h.TitleID != "" {
		fmt.Fprintf(&msg, "AJ%s%c", h.TitleID, delimiter)
	}

	if h.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", h.TerminalPassword, delimiter)
	}

	if h.FeeAcknowledged {
		fmt.Fprintf(&msg, "BO%s%c", utils.YorN(h.FeeAcknowledged), delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (h *Hold) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 27 {
		return 0, ErrInvalidRequest15
	}

	if string(runes[0:2]) != MsgIDHold {
		return 0, ErrInvalidRequest15
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "BW": "", "BS": "", "BY": "", "AO": "", "AA": "", "AD": "", "AB": "", "AJ": "", "AC": "", "BO": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	h.HoldMode = string(runes[2])

	h.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	if codes["BW"] != "" {
		h.ExpirationDate, err = time.Parse(utils.SIPDateFormat, codes["BW"])
		if err != nil {
			return 0, err
		}
	}

	h.PickupLocation = codes["BS"]
	if codes["BY"] != "" {
		h.HoldType, err = strconv.Atoi(codes["BY"])
		if err != nil {
			return 0, err
		}
	}

	h.InstitutionID = codes["AO"]
	h.PatronID = codes["AA"]

	h.PatronPassword = codes["AD"]
	h.ItemID = codes["AB"]
	h.TitleID = codes["AJ"]
	h.TerminalPassword = codes["AC"]

	if utf8.RuneCountInString(codes["BO"]) > 0 {
		h.FeeAcknowledged = utils.ParseBool([]rune(codes["BO"])[0])
	}

	err = h.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (h *Hold) Validate() error {
	err := Validate.Struct(h)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDHold, err.(validator.ValidationErrors))
	}
	return nil
}

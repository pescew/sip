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

var ErrInvalidResponse16 = fmt.Errorf("Invalid SIP %s", types.RespHold.String())

// The ACS should send this message in response to the Hold message from the SC.
type Hold struct {
	// Required Fields:
	Ok              bool
	Available       bool
	TransactionDate time.Time `validate:"required"`

	// Optional Fields:
	ExpirationDate time.Time
	QueuePosition  int    `validate:"min=-1"`
	PickupLocation string `validate:"sip"`
	InstitutionID  string `validate:"sip"`
	PatronID       string `validate:"sip"`
	ItemID         string `validate:"sip"`
	TitleID        string `validate:"sip"`
	ScreenMessage  string `validate:"sip"`
	PrintLine      string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (h *Hold) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder

	msg.WriteString(types.RespHold.ID())

	msg.WriteString(utils.ZeroOrOne(h.Ok))
	msg.WriteString(utils.YorN(h.Available))
	msg.WriteString(h.TransactionDate.Format(utils.SIPDateFormat))

	if !h.ExpirationDate.IsZero() {
		fmt.Fprintf(&msg, "BW%s%c", h.ExpirationDate.Format(utils.SIPDateFormat), delimiter)
	}

	if h.QueuePosition != -1 {
		fmt.Fprintf(&msg, "BR%d%c", h.QueuePosition, delimiter)
	}

	if h.PickupLocation != "" {
		fmt.Fprintf(&msg, "BS%s%c", h.PickupLocation, delimiter)
	}

	if h.InstitutionID != "" {
		fmt.Fprintf(&msg, "AO%s%c", h.InstitutionID, delimiter)
	}

	if h.PatronID != "" {
		fmt.Fprintf(&msg, "AA%s%c", h.PatronID, delimiter)
	}

	if h.ItemID != "" {
		fmt.Fprintf(&msg, "AB%s%c", h.ItemID, delimiter)
	}

	if h.TitleID != "" {
		fmt.Fprintf(&msg, "AJ%s%c", h.TitleID, delimiter)
	}

	if h.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", h.ScreenMessage, delimiter)
	}

	if h.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", h.PrintLine, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", h.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (h *Hold) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 22 {
		return ErrInvalidResponse16
	}

	if string(runes[0:2]) != types.RespHold.ID() {
		return ErrInvalidResponse16
	}

	var codes map[string]string
	if len(runes) > 22 {
		codes = utils.ExtractFields(string(runes[22:]), delimiter, map[string]string{"AY": "", "BW": "", "BR": "", "BS": "", "AO": "", "AA": "", "AB": "", "AJ": "", "AF": "", "AG": ""})
	} else {
		codes = map[string]string{"AY": "", "BW": "", "BR": "", "BS": "", "AO": "", "AA": "", "AB": "", "AJ": "", "AF": "", "AG": ""}
	}

	seqNumString := codes["AY"]
	if seqNumString == "" {
		h.SeqNum = 0
	} else {
		h.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			h.SeqNum = 0
		}
	}

	h.Ok = utils.ParseBool(runes[2])
	h.Available = utils.ParseBool(runes[3])

	h.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[4:22]))
	if err != nil {
		return err
	}

	if codes["BW"] != "" {
		h.ExpirationDate, err = time.Parse(utils.SIPDateFormat, codes["BW"])
		if err != nil {
			return err
		}
	}

	if codes["BR"] != "" {
		h.QueuePosition, err = strconv.Atoi(codes["BR"])
		if err != nil {
			return err
		}
	} else {
		h.QueuePosition = -1
	}

	h.PickupLocation = codes["BS"]
	h.InstitutionID = codes["AO"]
	h.PatronID = codes["AA"]
	h.ItemID = codes["AB"]
	h.TitleID = codes["AJ"]
	h.ScreenMessage = codes["AF"]
	h.PrintLine = codes["AG"]

	err = h.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (h *Hold) Validate() error {
	err := Validate.Struct(h)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespHold.String(), err.(validator.ValidationErrors))
	}
	return nil
}

package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDBlockPatron = "01"

var ErrInvalidRequest01 = fmt.Errorf("Invalid SIP %s request", MsgIDBlockPatron)

// This message requests that the patron card be blocked by the ACS. This is, for example, sent when the patron is detected tampering with the SC or when a patron forgets to take their card. The ACS should invalidate the patron’s card and respond with a Patron Status Response message. The ACS could also notify the library staff that the card has been blocked.
type BlockPatron struct {
	// Required:
	CardRetained     bool
	TransactionDate  time.Time `validate:"required"`
	InstitutionID    string    `validate:"required,sip"`
	BlockedCardMsg   string    `validate:"sip"`
	PatronID         string    `validate:"required,sip"`
	TerminalPassword string    `validate:"sip"`
}

func (bp *BlockPatron) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDBlockPatron)

	msg.WriteString(utils.YorN(bp.CardRetained))

	msg.WriteString(bp.TransactionDate.Format(utils.SIPDateFormat))

	fmt.Fprintf(&msg, "AO%s%c", bp.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AL%s%c", bp.BlockedCardMsg, delimiter)

	fmt.Fprintf(&msg, "AA%s%c", bp.PatronID, delimiter)
	fmt.Fprintf(&msg, "AC%s%c", bp.TerminalPassword, delimiter)

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (bp *BlockPatron) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 33 {
		return 0, ErrInvalidRequest01
	}

	if string(runes[0:2]) != MsgIDBlockPatron {
		return 0, ErrInvalidRequest01
	}

	codes := utils.ExtractFields(string(runes[21:]), delimiter, map[string]string{"AY": "", "AO": "", "AL": "", "AA": "", "AC": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	bp.CardRetained = utils.ParseBool(runes[2])

	bp.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[3:21]))
	if err != nil {
		return 0, err
	}

	bp.InstitutionID = codes["AO"]
	bp.BlockedCardMsg = codes["AL"]
	bp.PatronID = codes["AA"]
	bp.TerminalPassword = codes["AC"]

	err = bp.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (bp *BlockPatron) Validate() error {
	err := Validate.Struct(bp)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDBlockPatron, err.(validator.ValidationErrors))
	}
	return nil
}

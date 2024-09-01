package request

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

const MsgIDPatronInfo = "63"

var ErrInvalidRequest63 = fmt.Errorf("Invalid SIP %s request", MsgIDPatronInfo)

// This message is a superset of the Patron Status Request message. It should be used to request patron information. The ACS should respond with the Patron Information Response message.
type PatronInfo struct {
	// Required:
	Language        int            `validate:"min=0,max=999"`
	TransactionDate time.Time      `validate:"required"`
	Summary         fields.Summary `validate:"required"`
	InstitutionID   string         `validate:"required,sip"`
	PatronID        string         `validate:"required,sip"`

	// Optional:
	TerminalPassword string `validate:"sip"`
	PatronPassword   string `validate:"sip"`
	StartItem        int    `validate:"min=0"`
	EndItem          int    `validate:"min=0"`
}

func (pi *PatronInfo) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDPatronInfo)

	fmt.Fprintf(&msg, "%03d", pi.Language)
	msg.WriteString(pi.TransactionDate.Format(utils.SIPDateFormat))
	msg.WriteString(pi.Summary.Marshal())

	fmt.Fprintf(&msg, "AO%s%c", pi.InstitutionID, delimiter)
	fmt.Fprintf(&msg, "AA%s%c", pi.PatronID, delimiter)

	if pi.TerminalPassword != "" {
		fmt.Fprintf(&msg, "AC%s%c", pi.TerminalPassword, delimiter)
	}

	if pi.PatronPassword != "" {
		fmt.Fprintf(&msg, "AD%s%c", pi.PatronPassword, delimiter)
	}

	if pi.StartItem > 0 && pi.EndItem >= pi.StartItem {
		fmt.Fprintf(&msg, "BP%d%c", pi.StartItem, delimiter)
		fmt.Fprintf(&msg, "BQ%d%c", pi.EndItem, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (pi *PatronInfo) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 36 {
		return 0, ErrInvalidRequest63
	}

	if string(runes[0:2]) != MsgIDPatronInfo {
		return 0, ErrInvalidRequest63
	}

	codes := utils.ExtractFields(string(runes[33:]), delimiter, map[string]string{"AY": "", "AO": "", "AA": "", "AC": "", "AD": "", "BP": "", "BQ": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	pi.Language, err = strconv.Atoi(string(runes[2:5]))
	if err != nil {
		return 0, err
	}

	pi.TransactionDate, err = time.Parse(utils.SIPDateFormat, string(runes[5:23]))
	if err != nil {
		return 0, err
	}

	pi.Summary = fields.Summary{}
	err = pi.Summary.Unmarshal(string(runes[23:29]))
	if err != nil {
		return 0, err
	}

	pi.InstitutionID = codes["AO"]
	pi.PatronID = codes["AA"]
	pi.TerminalPassword = codes["AC"]
	pi.PatronPassword = codes["AD"]

	if codes["BP"] != "" {
		pi.StartItem, err = strconv.Atoi(codes["BP"])
		if err != nil {
			return 0, err
		}
	}

	if codes["BQ"] != "" {
		pi.EndItem, err = strconv.Atoi(codes["BQ"])
		if err != nil {
			return 0, err
		}
	}

	err = pi.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (pi *PatronInfo) Validate() error {
	err := pi.Summary.Validate()
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDPatronInfo, err)
	}

	err = Validate.Struct(pi)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDPatronInfo, err.(validator.ValidationErrors))
	}

	return nil
}

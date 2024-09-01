package request

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDSCStatus = "99"

var ErrInvalidRequest99 = fmt.Errorf("Invalid SIP %s request", MsgIDSCStatus)

// The SC status message sends SC status to the ACS. It requires an ACS Status Response message reply from the ACS. This message will be the first message sent by the SC to the ACS once a connection has been established (exception: the Login Message may be sent first to login to an ACS server program). The ACS will respond with a message that establishes some of the rules to be followed by the SC and establishes some parameters needed for further communication.
type SCStatus struct {
	// Required:
	StatusCode      int    `validate:"min=0,max=2"`
	MaxPrintWidth   int    `validate:"min=0,max=999"`
	ProtocolVersion string `validate:"required,sip,len=4,oneof=2.00"`
}

func (scs *SCStatus) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder
	msg.WriteString(MsgIDSCStatus)

	msg.WriteString(strconv.Itoa(scs.StatusCode))
	fmt.Fprintf(&msg, "%03d", scs.MaxPrintWidth)
	msg.WriteString(scs.ProtocolVersion)

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (scs *SCStatus) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 10 {
		return 0, ErrInvalidRequest99
	}

	if string(runes[0:2]) != MsgIDSCStatus {
		return 0, ErrInvalidRequest99
	}

	seqNum, err = strconv.Atoi(string(runes[12]))
	if err != nil {
		seqNum = 0
	}

	scs.StatusCode, err = strconv.Atoi(string(runes[2]))
	if err != nil {
		return 0, err
	}
	scs.MaxPrintWidth, err = strconv.Atoi(string(runes[3:6]))
	if err != nil {
		return 0, err
	}
	scs.ProtocolVersion = string(runes[6:10])

	err = scs.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (scs *SCStatus) Validate() error {
	err := Validate.Struct(scs)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDSCStatus, err.(validator.ValidationErrors))
	}
	return nil
}

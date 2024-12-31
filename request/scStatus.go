package request

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest99 = fmt.Errorf("Invalid SIP %s request", fields.ReqSCStatus.String())

// The SC status message sends SC status to the ACS. It requires an ACS Status Response message reply from the ACS. This message will be the first message sent by the SC to the ACS once a connection has been established (exception: the Login Message may be sent first to login to an ACS server program). The ACS will respond with a message that establishes some of the rules to be followed by the SC and establishes some parameters needed for further communication.
type SCStatus struct {
	// Required:
	StatusCode      fields.SCStatusCode    `validate:"min=0,max=2"`
	MaxPrintWidth   int                    `validate:"min=0,max=999"`
	ProtocolVersion fields.ProtocolVersion `validate:"required,sip,len=4,oneof=1.00 2.00"`

	SeqNum int `validate:"min=0,max=9"`
}

func (scs *SCStatus) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(fields.ReqSCStatus.ID())

	msg.WriteString(scs.StatusCode.ID())
	fmt.Fprintf(&msg, "%03d", scs.MaxPrintWidth)
	msg.WriteString(scs.ProtocolVersion.ID())

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", scs.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (scs *SCStatus) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 10 {
		return ErrInvalidRequest99
	}

	if string(runes[0:2]) != fields.ReqSCStatus.ID() {
		return ErrInvalidRequest99
	}

	scs.SeqNum, err = strconv.Atoi(string(runes[12]))
	if err != nil {
		scs.SeqNum = 0
	}

	scStatuscode, err := strconv.Atoi(string(runes[2]))
	if err != nil {
		return err
	}
	scs.StatusCode = fields.SCStatusCode(scStatuscode)

	scs.MaxPrintWidth, err = strconv.Atoi(string(runes[3:6]))
	if err != nil {
		return err
	}
	scs.ProtocolVersion = fields.ProtocolVersion(string(runes[6:10]))

	err = scs.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (scs *SCStatus) Validate() error {
	err := Validate.Struct(scs)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", fields.ReqSCStatus.String(), err.(validator.ValidationErrors))
	}
	return nil
}

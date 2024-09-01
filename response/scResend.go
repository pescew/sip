package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDSCResend = "96"

var ErrInvalidResponse96 = fmt.Errorf("Invalid SIP %s response", MsgIDSCResend)

// This message requests the SC to re-transmit its last message. It is sent by the ACS to the SC when the checksum in a received message does not match the value calculated by the ACS. The SC should respond by re-transmitting its last message, This message should never include a “sequence number” field, even when error detection is enabled, (see “Checksums and Sequence Numbers” below) but would include a “checksum” field since checksums are in use.
type SCResend struct{}

func (scr *SCResend) Marshal(seqNum int, delimiter, terminator rune) string {
	msg := MsgIDSCResend

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAZ", msg)), terminator)
}

func (scr *SCResend) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 2 {
		return 0, ErrInvalidResponse96
	}

	if string(runes[0:2]) != MsgIDSCResend {
		return 0, ErrInvalidResponse96
	}

	return seqNum, nil
}

func (scr *SCResend) Validate() error {
	err := Validate.Struct(scr)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDSCResend, err.(validator.ValidationErrors))
	}
	return nil
}

package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDACSResend = "97"

var ErrInvalidRequest97 = fmt.Errorf("Invalid SIP %s request", MsgIDACSResend)

// This message requests the ACS to re-transmit its last message. It is sent by the SC to the ACS when the checksum in a received message does not match the value calculated by the SC. The ACS should respond by re-transmitting its last message, This message should never include a “sequence number” field, even when error detection is enabled, (see “Checksums and Sequence Numbers” below) but would include a “checksum” field since checksums are in use.
type ACSResend struct{}

func (ar *ACSResend) Marshal(seqNum int, delimiter, terminator rune) string {
	msg := MsgIDACSResend

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAZ", msg)), terminator)
}

func (ar *ACSResend) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 2 {
		return 0, ErrInvalidRequest97
	}

	if string(runes[0:2]) != MsgIDACSResend {
		return 0, ErrInvalidRequest97
	}

	return seqNum, nil
}

func (ar *ACSResend) Validate() error {
	err := Validate.Struct(ar)
	if err != nil {
		return fmt.Errorf("invalid SIP %s request did not pass validation: %v", MsgIDACSResend, err.(validator.ValidationErrors))
	}
	return nil
}

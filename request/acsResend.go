package request

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest97 = fmt.Errorf("Invalid SIP %s request", types.ReqACSResend.String())

// This message requests the ACS to re-transmit its last message. It is sent by the SC to the ACS when the checksum in a received message does not match the value calculated by the SC. The ACS should respond by re-transmitting its last message, This message should never include a “sequence number” field, even when error detection is enabled, (see “Checksums and Sequence Numbers” below) but would include a “checksum” field since checksums are in use.
type ACSResend struct{}

func (ar *ACSResend) Marshal(delimiter, terminator rune, errorDetection bool) string {
	if errorDetection {
		var msg strings.Builder
		fmt.Fprintf(&msg, "%sAZ", types.ReqACSResend.ID())
		msg.WriteString(utils.ComputeChecksum(msg.String()))
		msg.WriteRune(terminator)
		return msg.String()
	}
	return fmt.Sprintf("%s%c", types.ReqACSResend.ID(), terminator)
}

func (ar *ACSResend) Unmarshal(line string, delimiter, terminator rune) error {
	runes := []rune(line)

	if len(runes) < 2 {
		return ErrInvalidRequest97
	}

	if string(runes[0:2]) != types.ReqACSResend.ID() {
		return ErrInvalidRequest97
	}

	return nil
}

func (ar *ACSResend) Validate() error {
	err := Validate.Struct(ar)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqACSResend.String(), err.(validator.ValidationErrors))
	}
	return nil
}

package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse96 = fmt.Errorf("Invalid SIP %s response", types.RespSCResend.String())

// This message requests the SC to re-transmit its last message. It is sent by the ACS to the SC when the checksum in a received message does not match the value calculated by the ACS. The SC should respond by re-transmitting its last message, This message should never include a “sequence number” field, even when error detection is enabled, (see “Checksums and Sequence Numbers” below) but would include a “checksum” field since checksums are in use.
type SCResend struct{}

func (scr *SCResend) Marshal(delimiter, terminator rune, errorDetection bool) string {
	if errorDetection {
		var msg strings.Builder
		fmt.Fprintf(&msg, "%sAZ", types.RespSCResend.ID())
		msg.WriteString(utils.ComputeChecksum(msg.String()))
		msg.WriteRune(terminator)
		return msg.String()
	}
	return fmt.Sprintf("%s%c", types.RespSCResend.ID(), terminator)
}

func (scr *SCResend) Unmarshal(line string, delimiter, terminator rune) error {
	runes := []rune(line)

	if len(runes) < 2 {
		return ErrInvalidResponse96
	}

	if string(runes[0:2]) != types.RespSCResend.ID() {
		return ErrInvalidResponse96
	}

	return nil
}

func (scr *SCResend) Validate() error {
	err := Validate.Struct(scr)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespSCResend.String(), err.(validator.ValidationErrors))
	}
	return nil
}

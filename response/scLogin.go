package response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidResponse94 = fmt.Errorf("Invalid SIP %s response", types.RespSCLogin.String())

// The ACS should send this message in response to the Login message. When this message is used, it will be the first message sent to the SC.
type SCLogin struct {
	// Required:
	Ok bool

	SeqNum int `validate:"min=0,max=9"`
}

func (scl *SCLogin) Marshal(delimiter, terminator rune, errorDetection bool) string {
	if errorDetection {
		var msg strings.Builder
		fmt.Fprintf(&msg, "%s%sAY%dAZ", types.RespSCLogin.ID(), utils.ZeroOrOne(scl.Ok), scl.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
		msg.WriteRune(terminator)
		return msg.String()
	}
	return fmt.Sprintf("%s%s%c", types.RespSCLogin.ID(), utils.ZeroOrOne(scl.Ok), terminator)
}

func (scl *SCLogin) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 3 {
		return ErrInvalidResponse94
	}

	if string(runes[0:2]) != types.RespSCLogin.ID() {
		return ErrInvalidResponse94
	}

	var codes map[string]string
	if len(runes) > 3 {
		codes = utils.ExtractFields(string(runes[3:]), delimiter, map[string]string{"AY": ""})
	} else {
		codes = map[string]string{"AY": ""}
	}

	seqNumString := codes["AY"]
	if seqNumString == "" {
		scl.SeqNum = 0
	} else {
		scl.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			scl.SeqNum = 0
		}
	}

	scl.Ok = utils.ParseBool(runes[2])

	err = scl.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (scl *SCLogin) Validate() error {
	err := Validate.Struct(scl)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.RespSCLogin.String(), err.(validator.ValidationErrors))
	}
	return nil
}

package response

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

const MsgIDSCLogin = "94"

var ErrInvalidResponse94 = fmt.Errorf("Invalid SIP %s response", MsgIDSCLogin)

// The ACS should send this message in response to the Login message. When this message is used, it will be the first message sent to the SC.
type SCLogin struct {
	// Required:
	Ok bool
}

func (scl *SCLogin) Marshal(seqNum int, delimiter, terminator rune) string {
	msg := MsgIDSCLogin

	msg += utils.ZeroOrOne(scl.Ok)

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg, seqNum)), terminator)
}

func (scl *SCLogin) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 3 {
		return 0, ErrInvalidResponse94
	}

	if string(runes[0:2]) != MsgIDSCLogin {
		return 0, ErrInvalidResponse94
	}

	var codes map[string]string
	if len(runes) > 3 {
		codes = utils.ExtractFields(string(runes[3:]), delimiter, map[string]string{"AY": ""})
	} else {
		codes = map[string]string{"AY": ""}
	}

	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	scl.Ok = utils.ParseBool(runes[2])

	err = scl.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (scl *SCLogin) Validate() error {
	err := Validate.Struct(scl)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDSCLogin, err.(validator.ValidationErrors))
	}
	return nil
}

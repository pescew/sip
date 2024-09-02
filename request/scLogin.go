package request

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var ErrInvalidRequest93 = fmt.Errorf("Invalid SIP %s request", types.ReqSCLogin.String())

// This message can be used to login to an ACS server program. The ACS should respond with the Login Response message. Whether to use this message or to use some other mechanism to login to the ACS is configurable on the SC. When this message is used, it will be the first message sent to the ACS.
type SCLogin struct {
	// Required:
	AlgorithmUserID   int    `validate:"min=0,max=9"`
	AlgorithmPassword int    `validate:"min=0,max=9"`
	LoginUserID       string `validate:"required,sip"`
	LoginPassword     string `validate:"sip"`

	// Optional:
	LocationCode string `validate:"sip"`

	SeqNum int `validate:"min=0,max=9"`
}

func (scl *SCLogin) Marshal(delimiter, terminator rune, errorDetection bool) string {
	var msg strings.Builder
	msg.WriteString(types.ReqSCLogin.ID())

	msg.WriteString(strconv.Itoa(scl.AlgorithmUserID))
	msg.WriteString(strconv.Itoa(scl.AlgorithmPassword))

	fmt.Fprintf(&msg, "CN%s%c", scl.LoginUserID, delimiter)
	fmt.Fprintf(&msg, "CO%s%c", scl.LoginPassword, delimiter)

	if scl.LocationCode != "" {
		fmt.Fprintf(&msg, "CP%s%c", scl.LocationCode, delimiter)
	}

	if errorDetection {
		fmt.Fprintf(&msg, "AY%dAZ", scl.SeqNum)
		msg.WriteString(utils.ComputeChecksum(msg.String()))
	}
	msg.WriteRune(terminator)
	return msg.String()
}

func (scl *SCLogin) Unmarshal(line string, delimiter, terminator rune) error {
	var err error
	runes := []rune(line)

	if len(runes) < 5 {
		return ErrInvalidRequest93
	}

	if string(runes[0:2]) != types.ReqSCLogin.ID() {
		return ErrInvalidRequest93
	}

	codes := utils.ExtractFields(string(runes[4:]), delimiter, map[string]string{"AY": "", "CN": "", "CO": "", "CP": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		scl.SeqNum = 0
	} else {
		scl.SeqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			scl.SeqNum = 0
		}
	}

	scl.AlgorithmUserID, err = strconv.Atoi(string(runes[2]))
	if err != nil {
		return err
	}
	scl.AlgorithmPassword, err = strconv.Atoi(string(runes[3]))
	if err != nil {
		return err
	}
	scl.LoginUserID = codes["CN"]
	scl.LoginPassword = codes["CO"]
	scl.LocationCode = codes["CP"]

	err = scl.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (scl *SCLogin) Validate() error {
	err := Validate.Struct(scl)
	if err != nil {
		return fmt.Errorf("invalid SIP %s did not pass validation: %v", types.ReqSCLogin.String(), err.(validator.ValidationErrors))
	}
	return nil
}

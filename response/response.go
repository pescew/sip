package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

var (
	Validate *validator.Validate

	ErrInvalidResponse = fmt.Errorf("Invalid SIP response")
	ErrUnknownResponse = fmt.Errorf("Unknown SIP response")
)

type Response interface {
	Marshal(delimiter, terminator rune, errorDetection bool) string
	Unmarshal(line string, delimiter, terminator rune) error
	Validate() error
}

func Unmarshal(line string, delimiter, terminator rune) (resp Response, msgID string, err error) {
	msgID = line[0:2]

	switch msgID {
	case types.RespCheckin.ID():
		resp = &Checkin{}
	case types.RespCheckout.ID():
		resp = &Checkout{}
	case types.RespHold.ID():
		resp = &Hold{}
	case types.RespItemInfo.ID():
		resp = &ItemInfo{}
	case types.RespItemStatusUpdate.ID():
		resp = &ItemStatusUpdate{}
	case types.RespPatronStatus.ID():
		resp = &PatronStatus{}
	case types.RespPatronEnable.ID():
		resp = &PatronEnable{}
	case types.RespRenew.ID():
		resp = &Renew{}
	case types.RespEndSession.ID():
		resp = &EndSession{}
	case types.RespFeePaid.ID():
		resp = &FeePaid{}
	case types.RespPatronInfo.ID():
		resp = &PatronInfo{}
	case types.RespRenewAll.ID():
		resp = &RenewAll{}
	case types.RespSCLogin.ID():
		resp = &SCLogin{}
	case types.RespSCResend.ID():
		resp = &SCResend{}
	case types.RespACSStatus.ID():
		resp = &ACSStatus{}
	default:
		return nil, msgID, ErrUnknownResponse
	}

	err = resp.Unmarshal(line, delimiter, terminator)
	if err != nil {
		return nil, msgID, err
	}

	return resp, msgID, nil
}

func InitValidator(excludeChars ...rune) {
	badChars := ""
	for _, char := range excludeChars {
		badChars += string(char)
	}

	Validate = validator.New()
	Validate.RegisterValidation("sip", utils.GenerateSIPValidatorFunc(badChars))
}

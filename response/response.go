package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
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
	case fields.RespCheckin.ID():
		resp = &Checkin{}
	case fields.RespCheckout.ID():
		resp = &Checkout{}
	case fields.RespHold.ID():
		resp = &Hold{}
	case fields.RespItemInfo.ID():
		resp = &ItemInfo{}
	case fields.RespItemStatusUpdate.ID():
		resp = &ItemStatusUpdate{}
	case fields.RespPatronStatus.ID():
		resp = &PatronStatus{}
	case fields.RespPatronEnable.ID():
		resp = &PatronEnable{}
	case fields.RespRenew.ID():
		resp = &Renew{}
	case fields.RespEndSession.ID():
		resp = &EndSession{}
	case fields.RespFeePaid.ID():
		resp = &FeePaid{}
	case fields.RespPatronInfo.ID():
		resp = &PatronInfo{}
	case fields.RespRenewAll.ID():
		resp = &RenewAll{}
	case fields.RespSCLogin.ID():
		resp = &SCLogin{}
	case fields.RespSCResend.ID():
		resp = &SCResend{}
	case fields.RespACSStatus.ID():
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

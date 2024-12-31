package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

var (
	Validate *validator.Validate

	ErrInvalidRequest = fmt.Errorf("Invalid SIP request")
	ErrUnknownRequest = fmt.Errorf("Unknown SIP request")
)

type Request interface {
	Marshal(delimiter, terminator rune, errorDetection bool) string
	Unmarshal(line string, delimiter, terminator rune) error
	Validate() error
}

func Unmarshal(line string, delimiter, terminator rune) (req Request, msgID string, err error) {
	msgID = line[0:2]

	switch msgID {
	case fields.ReqBlockPatron.ID():
		req = &BlockPatron{}
	case fields.ReqCheckin.ID():
		req = &Checkin{}
	case fields.ReqCheckout.ID():
		req = &Checkout{}
	case fields.ReqHold.ID():
		req = &Hold{}
	case fields.ReqItemInfo.ID():
		req = &ItemInfo{}
	case fields.ReqItemStatusUpdate.ID():
		req = &ItemStatusUpdate{}
	case fields.ReqPatronStatus.ID():
		req = &PatronStatus{}
	case fields.ReqPatronEnable.ID():
		req = &PatronEnable{}
	case fields.ReqRenew.ID():
		req = &Renew{}
	case fields.ReqEndPatronSession.ID():
		req = &EndPatronSession{}
	case fields.ReqFeePaid.ID():
		req = &FeePaid{}
	case fields.ReqPatronInfo.ID():
		req = &PatronInfo{}
	case fields.ReqRenewAll.ID():
		req = &RenewAll{}
	case fields.ReqSCLogin.ID():
		req = &SCLogin{}
	case fields.ReqACSResend.ID():
		req = &ACSResend{}
	case fields.ReqSCStatus.ID():
		req = &SCStatus{}
	default:
		return nil, msgID, ErrUnknownRequest
	}

	err = req.Unmarshal(line, delimiter, terminator)
	if err != nil {
		return nil, msgID, err
	}

	return req, msgID, nil
}

func InitValidator(excludeChars ...rune) {
	badChars := ""
	for _, char := range excludeChars {
		badChars += string(char)
	}

	Validate = validator.New()
	Validate.RegisterValidation("sip", utils.GenerateSIPValidatorFunc(badChars))
}

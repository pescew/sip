package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/types"
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
	case types.ReqBlockPatron.ID():
		req = &BlockPatron{}
	case types.ReqCheckin.ID():
		req = &Checkin{}
	case types.ReqCheckout.ID():
		req = &Checkout{}
	case types.ReqHold.ID():
		req = &Hold{}
	case types.ReqItemInfo.ID():
		req = &ItemInfo{}
	case types.ReqItemStatusUpdate.ID():
		req = &ItemStatusUpdate{}
	case types.ReqPatronStatus.ID():
		req = &PatronStatus{}
	case types.ReqPatronEnable.ID():
		req = &PatronEnable{}
	case types.ReqRenew.ID():
		req = &Renew{}
	case types.ReqEndPatronSession.ID():
		req = &EndPatronSession{}
	case types.ReqFeePaid.ID():
		req = &FeePaid{}
	case types.ReqPatronInfo.ID():
		req = &PatronInfo{}
	case types.ReqRenewAll.ID():
		req = &RenewAll{}
	case types.ReqSCLogin.ID():
		req = &SCLogin{}
	case types.ReqACSResend.ID():
		req = &ACSResend{}
	case types.ReqSCStatus.ID():
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

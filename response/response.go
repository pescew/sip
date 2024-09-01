package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

var (
	Validate *validator.Validate

	ErrInvalidResponse = fmt.Errorf("Invalid SIP response")
	ErrUnknownResponse = fmt.Errorf("Unknown SIP response")
)

type Response interface {
	Marshal(seqNum int, delimiter, terminator rune) string
	Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error)
	Validate() error
}

func Unmarshal(line string, delimiter, terminator rune) (resp Response, msgID string, seqNum int, err error) {
	msgID = line[0:2]

	switch msgID {
	case MsgIDCheckin:
		resp = &Checkin{}
	case MsgIDCheckout:
		resp = &Checkout{}
	case MsgIDHold:
		resp = &Hold{}
	case MsgIDItemInfo:
		resp = &ItemInfo{}
	case MsgIDItemStatusUpdate:
		resp = &ItemStatusUpdate{}
	case MsgIDPatronStatus:
		resp = &PatronStatus{}
	case MsgIDPatronEnable:
		resp = &PatronEnable{}
	case MsgIDRenew:
		resp = &Renew{}
	case MsgIDEndSession:
		resp = &EndSession{}
	case MsgIDFeePaid:
		resp = &FeePaid{}
	case MsgIDPatronInfo:
		resp = &PatronInfo{}
	case MsgIDRenewAll:
		resp = &RenewAll{}
	case MsgIDSCLogin:
		resp = &SCLogin{}
	case MsgIDSCResend:
		resp = &SCResend{}
	case MsgIDACSStatus:
		resp = &ACSStatus{}
	default:
		return nil, "", 0, ErrUnknownResponse
	}

	seqNum, err = resp.Unmarshal(line, delimiter, terminator)
	if err != nil {
		return nil, "", 0, err
	}

	return resp, msgID, seqNum, nil
}

func InitValidator(excludeChars ...rune) {
	badChars := ""
	for _, char := range excludeChars {
		badChars += string(char)
	}

	Validate = validator.New()
	Validate.RegisterValidation("sip", utils.GenerateSIPValidatorFunc(badChars))
}

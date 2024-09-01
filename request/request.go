package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/utils"
)

var (
	Validate *validator.Validate

	ErrInvalidRequest = fmt.Errorf("Invalid SIP request")
	ErrUnknownRequest = fmt.Errorf("Unknown SIP request")
)

type Request interface {
	Marshal(seqNum int, delimiter, terminator rune) string
	Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error)
	Validate() error
}

func Unmarshal(line string, delimiter, terminator rune) (req Request, msgID string, seqNum int, err error) {
	msgID = line[0:2]

	switch msgID {
	case MsgIDBlockPatron:
		req = &BlockPatron{}
	case MsgIDCheckin:
		req = &Checkin{}
	case MsgIDCheckout:
		req = &Checkout{}
	case MsgIDHold:
		req = &Hold{}
	case MsgIDItemInfo:
		req = &ItemInfo{}
	case MsgIDItemStatusUpdate:
		req = &ItemStatusUpdate{}
	case MsgIDPatronStatus:
		req = &PatronStatus{}
	case MsgIDPatronEnable:
		req = &PatronEnable{}
	case MsgIDRenew:
		req = &Renew{}
	case MsgIDEndPatronSession:
		req = &EndPatronSession{}
	case MsgIDFeePaid:
		req = &FeePaid{}
	case MsgIDPatronInfo:
		req = &PatronInfo{}
	case MsgIDRenewAll:
		req = &RenewAll{}
	case MsgIDSCLogin:
		req = &SCLogin{}
	case MsgIDACSResend:
		req = &ACSResend{}
	case MsgIDSCStatus:
		req = &SCStatus{}
	default:
		return nil, "", 0, ErrUnknownRequest
	}

	seqNum, err = req.Unmarshal(line, delimiter, terminator)
	if err != nil {
		return nil, "", 0, err
	}

	return req, msgID, seqNum, nil
}

func InitValidator(excludeChars ...rune) {
	badChars := ""
	for _, char := range excludeChars {
		badChars += string(char)
	}

	Validate = validator.New()
	Validate.RegisterValidation("sip", utils.GenerateSIPValidatorFunc(badChars))
}

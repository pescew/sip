package fields

import (
	"fmt"
	"strings"

	"github.com/pescew/sip/utils"
)

var ErrInvalidFieldSupportedMessages = fmt.Errorf("Invalid SIP field: supportedMessages")

// Variable-length field. This field is used to notify the SC about which messages the ACS supports. A Y in a position means that the associated message/response is supported. An N means the message/response pair is not supported.
type SupportedMessages struct {
	PatronStatusRequest bool
	Checkout            bool
	Checkin             bool
	BlockPatron         bool
	SCACSStatus         bool
	RequestResend       bool
	Login               bool
	PatronInformation   bool
	EndPatronSession    bool
	FeePaid             bool
	ItemInformation     bool
	ItemStatusUpdate    bool
	PatronEnable        bool
	Hold                bool
	Renew               bool
	RenewAll            bool
}

func (sm *SupportedMessages) Marshal() string {
	var msg strings.Builder

	msg.WriteString(utils.YorN(sm.PatronStatusRequest))
	msg.WriteString(utils.YorN(sm.Checkout))
	msg.WriteString(utils.YorN(sm.Checkin))
	msg.WriteString(utils.YorN(sm.BlockPatron))
	msg.WriteString(utils.YorN(sm.SCACSStatus))
	msg.WriteString(utils.YorN(sm.RequestResend))
	msg.WriteString(utils.YorN(sm.Login))
	msg.WriteString(utils.YorN(sm.PatronInformation))
	msg.WriteString(utils.YorN(sm.EndPatronSession))
	msg.WriteString(utils.YorN(sm.FeePaid))
	msg.WriteString(utils.YorN(sm.ItemInformation))
	msg.WriteString(utils.YorN(sm.ItemStatusUpdate))
	msg.WriteString(utils.YorN(sm.PatronEnable))
	msg.WriteString(utils.YorN(sm.Hold))
	msg.WriteString(utils.YorN(sm.Renew))
	msg.WriteString(utils.YorN(sm.RenewAll))

	return msg.String()
}

func (sm *SupportedMessages) Unmarshal(line string) {
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		switch i {
		case 0:
			sm.PatronStatusRequest = utils.ParseBool(runes[0])
		case 1:
			sm.Checkout = utils.ParseBool(runes[1])
		case 2:
			sm.Checkin = utils.ParseBool(runes[2])
		case 3:
			sm.BlockPatron = utils.ParseBool(runes[3])
		case 4:
			sm.SCACSStatus = utils.ParseBool(runes[4])
		case 5:
			sm.RequestResend = utils.ParseBool(runes[5])
		case 6:
			sm.Login = utils.ParseBool(runes[6])
		case 7:
			sm.PatronInformation = utils.ParseBool(runes[7])
		case 8:
			sm.EndPatronSession = utils.ParseBool(runes[8])
		case 9:
			sm.FeePaid = utils.ParseBool(runes[9])
		case 10:
			sm.ItemInformation = utils.ParseBool(runes[10])
		case 11:
			sm.ItemStatusUpdate = utils.ParseBool(runes[11])
		case 12:
			sm.PatronEnable = utils.ParseBool(runes[12])
		case 13:
			sm.Hold = utils.ParseBool(runes[13])
		case 14:
			sm.Renew = utils.ParseBool(runes[14])
		case 15:
			sm.RenewAll = utils.ParseBool(runes[15])
		default:
			//do nothing
		}
	}
}

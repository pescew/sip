package fields

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pescew/sip/utils"
)

var ErrInvalidFieldPatronStatus = fmt.Errorf("Invalid SIP field: patronStatus")

// 14-char, fixed-length field. This field is described in the preliminary NISO standard Z39.70-199x. A Y in any position indicates that the condition is true. A blank (code $20) in this position means that this condition is not true. For example, the first position of this field corresponds to "charge privileges denied" and must therefore contain a code $20 if this patron’s privileges are authorized.
type PatronStatus struct {
	DenyCharges           bool
	DenyRenewals          bool
	DenyRecalls           bool
	DenyHolds             bool
	CardLost              bool
	TooManyCharged        bool
	TooManyOverdue        bool
	TooManyRenewals       bool
	TooManyClaimsReturned bool
	TooManyItemsLost      bool
	ExceedsFines          bool
	ExceedsFees           bool
	RecallOverdue         bool
	TooManyBilled         bool
}

func (ps *PatronStatus) Marshal() string {
	var msg strings.Builder

	msg.WriteString(utils.YorN(ps.DenyCharges))
	msg.WriteString(utils.YorN(ps.DenyRenewals))
	msg.WriteString(utils.YorN(ps.DenyRecalls))
	msg.WriteString(utils.YorN(ps.DenyHolds))
	msg.WriteString(utils.YorN(ps.CardLost))
	msg.WriteString(utils.YorN(ps.TooManyCharged))
	msg.WriteString(utils.YorN(ps.TooManyOverdue))
	msg.WriteString(utils.YorN(ps.TooManyRenewals))
	msg.WriteString(utils.YorN(ps.TooManyClaimsReturned))
	msg.WriteString(utils.YorN(ps.TooManyItemsLost))
	msg.WriteString(utils.YorN(ps.ExceedsFines))
	msg.WriteString(utils.YorN(ps.ExceedsFees))
	msg.WriteString(utils.YorN(ps.RecallOverdue))
	msg.WriteString(utils.YorN(ps.TooManyBilled))

	return msg.String()
}

func (ps *PatronStatus) Unmarshal(line string) error {
	if utf8.RuneCountInString(line) < 14 {
		return ErrInvalidFieldPatronStatus
	}

	runes := []rune(line)

	ps.DenyCharges = utils.ParseBool(runes[0])
	ps.DenyRenewals = utils.ParseBool(runes[1])
	ps.DenyRecalls = utils.ParseBool(runes[2])
	ps.DenyHolds = utils.ParseBool(runes[3])
	ps.CardLost = utils.ParseBool(runes[4])
	ps.TooManyCharged = utils.ParseBool(runes[5])
	ps.TooManyOverdue = utils.ParseBool(runes[6])
	ps.TooManyRenewals = utils.ParseBool(runes[7])
	ps.TooManyClaimsReturned = utils.ParseBool(runes[8])
	ps.TooManyItemsLost = utils.ParseBool(runes[9])
	ps.ExceedsFines = utils.ParseBool(runes[10])
	ps.ExceedsFees = utils.ParseBool(runes[11])
	ps.RecallOverdue = utils.ParseBool(runes[12])
	ps.TooManyBilled = utils.ParseBool(runes[13])

	return nil
}

package fields

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pescew/sip/utils"
)

var ErrInvalidFieldSummary = fmt.Errorf("Invalid SIP field: summary")

// 10-char, fixed-length field. This allows the SC to request partial information only. This field usage is similar to the NISO defined PATRON STATUS field. A Y in any position indicates that detailed as well as summary information about the corresponding category of items can be sent in the response. A blank (code $20) in this position means that only summary information should be sent about the corresponding category of items. Only one category of items should be requested at a time, i.e. it would take 6 of these messages, each with a different position set to Y, to get all the detailed information about a patron’s items. All of the 6 responses, however, would contain the summary information. See Patron Information Response.
type Summary struct {
	HoldItems        bool
	OverdueItems     bool
	ChargedItems     bool
	FineItems        bool
	RecallItems      bool
	UnavailableHolds bool
}

func (s *Summary) Marshal() string {
	var msg strings.Builder

	msg.WriteString(utils.YorN(s.HoldItems))
	msg.WriteString(utils.YorN(s.OverdueItems))
	msg.WriteString(utils.YorN(s.ChargedItems))
	msg.WriteString(utils.YorN(s.FineItems))
	msg.WriteString(utils.YorN(s.RecallItems))
	msg.WriteString(utils.YorN(s.UnavailableHolds))
	msg.WriteString("    ")

	return msg.String()
}

func (s *Summary) Unmarshal(line string) error {
	if utf8.RuneCountInString(line) < 6 {
		return ErrInvalidFieldSummary
	}

	runes := []rune(line)

	s.HoldItems = utils.ParseBool(runes[0])
	s.OverdueItems = utils.ParseBool(runes[1])
	s.ChargedItems = utils.ParseBool(runes[2])
	s.FineItems = utils.ParseBool(runes[3])
	s.RecallItems = utils.ParseBool(runes[4])
	s.UnavailableHolds = utils.ParseBool(runes[5])

	return nil
}

func (s *Summary) Validate() error {
	numTrue := 0
	if s.HoldItems {
		numTrue++
	}
	if s.OverdueItems {
		numTrue++
	}
	if s.ChargedItems {
		numTrue++
	}
	if s.FineItems {
		numTrue++
	}
	if s.RecallItems {
		numTrue++
	}
	if s.UnavailableHolds {
		numTrue++
	}
	if numTrue > 1 {
		return fmt.Errorf("invalid SIP summary field. Max 1 summary per request, %d requested. %v", numTrue, ErrInvalidFieldSummary)
	}

	return nil
}

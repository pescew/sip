package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

func TestACSStatus(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *ACSStatus
	resp := &ACSStatus{
		// Required:
		OnlineStatus:    true,
		CheckinOK:       true,
		CheckoutOK:      true,
		RenewalPolicy:   true,
		StatusUpdateOK:  false,
		OfflineOK:       false,
		TimeoutPeriod:   744,
		RetriesAllowed:  45,
		DateTimeSync:    time.Now().UTC().Truncate(time.Second),
		ProtocolVersion: "2.00",
		InstitutionID:   "inst",

		// Optional:
		LibraryName: "lib",

		// Required:
		SupportedMessages: fields.SupportedMessages{
			PatronStatusRequest: true,
			Checkout:            true,
			Checkin:             true,
			BlockPatron:         true,
			SCACSStatus:         true,
			RequestResend:       false,
			Login:               true,
			PatronInformation:   false,
			EndPatronSession:    true,
			FeePaid:             true,
			ItemInformation:     true,
			ItemStatusUpdate:    false,
			PatronEnable:        true,
			Hold:                true,
			Renew:               false,
			RenewAll:            true,
		},

		// Optional:
		TerminalLocation: "lib",
		ScreenMessage:    "",
		PrintLine:        "",
	}

	sipString := resp.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*ACSStatus)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDACSStatus {
		t.Fatalf("Message ID mismatch")
	}

	if !cmp.Equal(resp, respParsed) {
		fmt.Println("----------")
		fmt.Println(resp)
		fmt.Println("----------")
		fmt.Println(sipString)
		fmt.Println("----------")
		fmt.Println(respParsed)
		fmt.Println("----------")
		t.Fatalf("struct mismatch")
	}
}

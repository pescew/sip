package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestPatronInfo(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *PatronInfo
	req := &PatronInfo{
		// Required:
		Language:        0,
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		Summary: fields.Summary{
			HoldItems:        true,
			OverdueItems:     false,
			ChargedItems:     false,
			FineItems:        false,
			RecallItems:      false,
			UnavailableHolds: false,
		},
		InstitutionID: "inst",
		PatronID:      "johndoe",

		// Optional:
		TerminalPassword: "password",
		PatronPassword:   "john'sPassword",
		StartItem:        2,
		EndItem:          4,

		SeqNum: 3,
	}

	sipString := req.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*PatronInfo)

	if reqParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.ReqPatronInfo.ID() {
		t.Fatalf("Message ID mismatch")
	}

	if !cmp.Equal(req, reqParsed) {
		fmt.Println("----------")
		fmt.Println(req)
		fmt.Println("----------")
		fmt.Println(sipString)
		fmt.Println("----------")
		fmt.Println(reqParsed)
		fmt.Println("----------")
		t.Fatalf("struct mismatch")
	}
}

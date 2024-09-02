package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestHold(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *Hold
	req := &Hold{
		// Required:
		HoldMode:        "+",
		TransactionDate: time.Now().UTC().Truncate(time.Second),

		// Optional:
		ExpirationDate: time.Now().UTC().Truncate(time.Second),
		PickupLocation: "lib",
		HoldType:       3,

		// Required:
		InstitutionID: "inst",
		PatronID:      "johndoe",

		// Optional:
		PatronPassword:   "john'sPassword",
		ItemID:           "1234567890",
		TitleID:          "",
		TerminalPassword: "password",
		FeeAcknowledged:  true,

		SeqNum: 3,
	}

	sipString := req.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*Hold)

	if reqParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.ReqHold.ID() {
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

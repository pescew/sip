package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestHold(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *Hold
	resp := &Hold{
		// Required Fields:
		Ok:              true,
		Available:       true,
		TransactionDate: time.Now().UTC().Truncate(time.Second),

		// Optional Fields:
		ExpirationDate: time.Now().UTC().Truncate(time.Second),
		QueuePosition:  12,
		PickupLocation: "lib",
		InstitutionID:  "inst",
		PatronID:       "0987654321",
		ItemID:         "1234567890",
		TitleID:        "Item Title",
		ScreenMessage:  "msg",
		PrintLine:      "print",
	}

	sipString := resp.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*Hold)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDHold {
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

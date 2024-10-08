package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestFeePaid(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *FeePaid
	resp := &FeePaid{
		// Required Fields:
		PaymentAccepted: true,
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		InstitutionID:   "inst",
		PatronID:        "0987654321",

		// Optional Fields:
		TransactionID: "12345",
		ScreenMessage: "msg",
		PrintLine:     "print",

		SeqNum: 3,
	}

	sipString := resp.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*FeePaid)

	if respParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.RespFeePaid.ID() {
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

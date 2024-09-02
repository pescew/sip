package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestRenewAll(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *RenewAll
	resp := &RenewAll{
		// Required Fields:
		Ok:              true,
		RenewedCount:    3,
		UnrenewedCount:  0,
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		InstitutionID:   "inst",

		// Optional Fields:
		RenewedItems:   []string{"1234567890", "0987654321", "5555555555"},
		UnrenewedItems: []string{"0987654321", "5555555555"},
		ScreenMessage:  "msg",
		PrintLine:      "print",

		SeqNum: 3,
	}

	sipString := resp.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*RenewAll)

	if respParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.RespRenewAll.ID() {
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

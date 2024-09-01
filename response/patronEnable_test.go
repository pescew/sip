package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

func TestPatronEnable(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *PatronEnable
	resp := &PatronEnable{
		// Required Fields:
		PatronStatus: fields.PatronStatus{
			DenyCharges:           true,
			DenyRenewals:          true,
			DenyRecalls:           false,
			DenyHolds:             true,
			CardLost:              true,
			TooManyCharged:        true,
			TooManyOverdue:        false,
			TooManyRenewals:       true,
			TooManyClaimsReturned: true,
			TooManyItemsLost:      true,
			ExceedsFines:          false,
			ExceedsFees:           true,
			RecallOverdue:         true,
			TooManyBilled:         true,
		},
		Language:        1,
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		InstitutionID:   "inst",
		PatronID:        "0987654321",
		PatronName:      "Doe, John",

		// Optional Fields:
		ValidPatron:         true,
		ValidPatronPassword: true,
		ScreenMessage:       "msg",
		PrintLine:           "print",
	}

	sipString := resp.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*PatronEnable)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDPatronEnable {
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

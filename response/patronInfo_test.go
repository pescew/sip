package response

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
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *PatronInfo
	resp := &PatronInfo{
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
		Language:              1,
		TransactionDate:       time.Now().UTC().Truncate(time.Second),
		HoldItemsCount:        2,
		OverdueItemsCount:     0,
		ChargedItemsCount:     1,
		FineItemsCount:        1,
		RecallItemsCount:      0,
		UnavailableHoldsCount: 6,
		InstitutionID:         "inst",
		PatronID:              "0987654321",
		PatronName:            "Doe, John",

		// Optional Fields:
		HoldItemsLimit:      50,
		OverdueItemsLimit:   50,
		ChargedItemsLimit:   50,
		ValidPatron:         true,
		ValidPatronPassword: true,
		CurrencyType:        "USD",
		FeeAmount:           "25.50",
		FeeLimit:            "50.00",
		HoldItems:           []string{"1234567890", "0987654321", "5555555555"},
		OverdueItems:        []string{"0987654321", "5555555555"},
		ChargedItems:        []string{"1234567890"},
		FineItems:           []string{"1234567890"},
		RecallItems:         []string{},
		UnavailHoldItems:    []string{"1111111111", "2222222222", "3333333333", "4444444444"},
		HomeAddress:         "123 Main Street, New York, NY 11111",
		EmailAddress:        "test@test.com",
		HomePhone:           "555-555-5555",
		ScreenMessage:       "msg",
		PrintLine:           "print",

		SeqNum: 3,
	}

	sipString := resp.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*PatronInfo)

	if respParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.RespPatronInfo.ID() {
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

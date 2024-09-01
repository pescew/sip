package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestCheckout(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *Checkout
	resp := &Checkout{
		// Required Fields:
		Ok:              true,
		RenewalOk:       true,
		MagneticMedia:   true,
		Desensitize:     false,
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		InstitutionID:   "inst",
		PatronID:        "0987654321",
		ItemID:          "1234567890",
		TitleID:         "Item Title",
		DueDate:         "12/31/1969",

		// Optional Fields:
		FeeType:         12,
		SecurityInhibit: false,
		CurrencyType:    "USD",
		FeeAmount:       "50.00",
		MediaType:       "005",
		ItemProperties:  "props",
		TransactionID:   "12345",
		ScreenMessage:   "msg",
		PrintLine:       "print",
	}

	sipString := resp.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*Checkout)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDCheckout {
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

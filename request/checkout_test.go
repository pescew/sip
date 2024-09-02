package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestCheckout(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *Checkout
	req := &Checkout{
		// Required:
		SCRenewalPolicy:  true,
		NoBlock:          true,
		TransactionDate:  time.Now().UTC().Truncate(time.Second),
		NBDueDate:        time.Now().UTC().Truncate(time.Second),
		InstitutionID:    "inst",
		PatronID:         "johndoe",
		ItemID:           "1234567890",
		TerminalPassword: "password",

		// Optional:
		ItemProperties:  "",
		PatronPassword:  "john'sPassword",
		FeeAcknowledged: true,
		Cancel:          false,

		SeqNum: 3,
	}

	sipString := req.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*Checkout)

	if reqParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.ReqCheckout.ID() {
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

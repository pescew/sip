package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestFeePaid(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *FeePaid
	req := &FeePaid{
		// Required:
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		FeeType:         4,
		PaymentType:     2,
		CurrencyType:    "USD",
		FeeAmount:       "50.00",
		InstitutionID:   "inst",
		PatronID:        "johndoe",

		// Optional:
		TerminalPassword: "password",
		PatronPassword:   "john'sPassword",
		FeeID:            "523w44fghdf",
		TransactionID:    "sdgf345ydfhg6",

		SeqNum: 3,
	}

	sipString := req.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*FeePaid)

	if reqParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.ReqFeePaid.ID() {
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

package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestItemInfo(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *ItemInfo
	req := &ItemInfo{
		// Required:
		TransactionDate: time.Now().UTC().Truncate(time.Second),
		InstitutionID:   "inst",
		ItemID:          "1234567890",

		// Optional:
		TerminalPassword: "password",
	}

	sipString := req.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*ItemInfo)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDItemInfo {
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

package request

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestSCStatus(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *SCStatus
	req := &SCStatus{
		// Required:
		StatusCode:      1,
		MaxPrintWidth:   30,
		ProtocolVersion: "2.00",
	}

	sipString := req.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*SCStatus)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch: %d != %d", seqNum, 3)
	}

	if msgID != MsgIDSCStatus {
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

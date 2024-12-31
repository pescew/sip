package request

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/fields"
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
		StatusCode:      fields.SCStatusOK,
		MaxPrintWidth:   30,
		ProtocolVersion: fields.ProtocolVersion2,

		SeqNum: 3,
	}

	sipString := req.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*SCStatus)

	if reqParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch: %d != %d", reqParsed.SeqNum, 3)
	}

	if msgID != fields.ReqSCStatus.ID() {
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

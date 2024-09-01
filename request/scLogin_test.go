package request

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestSCLogin(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *SCLogin
	req := &SCLogin{
		// Required:
		AlgorithmUserID:   0,
		AlgorithmPassword: 0,
		LoginUserID:       "testUser",
		LoginPassword:     "testPass",

		// Optional:
		LocationCode: "lib",
	}

	sipString := req.Marshal(3, delimiter, terminator)

	parsed, msgID, seqNum, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*SCLogin)

	if seqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != MsgIDSCLogin {
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

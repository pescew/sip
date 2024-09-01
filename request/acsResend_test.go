package request

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestACSResend(t *testing.T) {
	delimiter := []rune("|")[0]
	terminator := []rune("\r")[0]

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var reqParsed *ACSResend
	req := &ACSResend{}

	sipString := req.Marshal(3, delimiter, terminator)

	parsed, msgID, _, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	reqParsed = parsed.(*ACSResend)

	if msgID != MsgIDACSResend {
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

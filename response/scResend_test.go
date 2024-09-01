package response

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/utils"
)

func TestSCResend(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *SCResend
	resp := &SCResend{}

	sipString := resp.Marshal(3, delimiter, terminator)

	parsed, msgID, _, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*SCResend)

	if msgID != MsgIDSCResend {
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

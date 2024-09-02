package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestCheckin(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *Checkin
	resp := &Checkin{
		// Required Fields:
		Ok:                true,
		Resensitize:       true,
		MagneticMedia:     true,
		Alert:             false,
		TransactionDate:   time.Now().UTC().Truncate(time.Second),
		InstitutionID:     "inst",
		ItemID:            "1234567890",
		PermanentLocation: "lib",

		// Optional Fields:
		TitleID:        "Item Title",
		SortBin:        "4",
		PatronID:       "0987654321",
		MediaType:      "005",
		ItemProperties: "props",
		ScreenMessage:  "msg",
		PrintLine:      "print",

		SeqNum: 3,
	}

	sipString := resp.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*Checkin)

	if respParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.RespCheckin.ID() {
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

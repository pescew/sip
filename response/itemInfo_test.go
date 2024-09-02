package response

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func TestItemInfo(t *testing.T) {
	delimiter := '|'
	terminator := '\r'

	InitValidator(delimiter, terminator)
	utils.ConfigureEscapeCharacters(delimiter, terminator)

	var respParsed *ItemInfo
	resp := &ItemInfo{
		// Required Fields:
		CirculationStatus: 28,
		SecurityMarker:    83,
		FeeType:           12,
		TransactionDate:   time.Now().UTC().Truncate(time.Second),

		// Optional Fields:
		HoldQueueLength: 5,
		DueDate:         "12/31/1969",
		RecallDate:      time.Now().UTC().Truncate(time.Second),
		HoldPickupDate:  time.Now().UTC().Truncate(time.Second),

		// Required Fields:
		ItemID:  "1234567890",
		TitleID: "Item Title",

		// Optional Fields:
		Owner:             "lib",
		CurrencyType:      "USD",
		FeeAmount:         "50.00",
		MediaType:         "005",
		PermanentLocation: "lib1",
		CurrentLocation:   "lib2",
		ItemProperties:    "props",
		ScreenMessage:     "msg",
		PrintLine:         "print",

		SeqNum: 3,
	}

	sipString := resp.Marshal(delimiter, terminator, true)

	parsed, msgID, err := Unmarshal(sipString, delimiter, terminator)
	if err != nil {
		t.Fatal(err)
	}

	respParsed = parsed.(*ItemInfo)

	if respParsed.SeqNum != 3 {
		t.Fatalf("Sequence Number mismatch")
	}

	if msgID != types.RespItemInfo.ID() {
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

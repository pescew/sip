package response

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/utils"
)

const MsgIDACSStatus = "98"

var ErrInvalidResponse98 = fmt.Errorf("Invalid SIP %s response", MsgIDACSStatus)

// The ACS must send this message in response to a SC Status message. This message will be the first message sent by the ACS to the SC, since it establishes some of the rules to be followed by the SC and establishes some parameters needed for further communication (exception: the Login Response Message may be sent first to complete login of the SC).
type ACSStatus struct {
	// Required:
	OnlineStatus    bool
	CheckinOK       bool
	CheckoutOK      bool
	RenewalPolicy   bool
	StatusUpdateOK  bool
	OfflineOK       bool
	TimeoutPeriod   int       `validate:"min=0,max=999"`
	RetriesAllowed  int       `validate:"min=0,max=999"`
	DateTimeSync    time.Time `validate:"required"`
	ProtocolVersion string    `validate:"required,sip,len=4,oneof=2.00"`
	InstitutionID   string    `validate:"required,sip"`

	// Optional:
	LibraryName string `validate:"sip"`

	// Required:
	SupportedMessages fields.SupportedMessages `validate:"required"`

	// Optional:
	TerminalLocation string `validate:"sip"`
	ScreenMessage    string `validate:"sip"`
	PrintLine        string `validate:"sip"`
}

func (st *ACSStatus) Marshal(seqNum int, delimiter, terminator rune) string {
	var msg strings.Builder

	msg.WriteString(MsgIDACSStatus)

	msg.WriteString(utils.YorN(st.OnlineStatus))
	msg.WriteString(utils.YorN(st.CheckinOK))
	msg.WriteString(utils.YorN(st.CheckoutOK))
	msg.WriteString(utils.YorN(st.RenewalPolicy))
	msg.WriteString(utils.YorN(st.StatusUpdateOK))
	msg.WriteString(utils.YorN(st.OfflineOK))

	fmt.Fprintf(&msg, "%03d", st.TimeoutPeriod)
	fmt.Fprintf(&msg, "%03d", st.RetriesAllowed)

	msg.WriteString(st.DateTimeSync.Format(utils.SIPDateFormat))

	msg.WriteString(st.ProtocolVersion)

	fmt.Fprintf(&msg, "AO%s%c", st.InstitutionID, delimiter)

	if st.LibraryName != "" {
		fmt.Fprintf(&msg, "AM%s%c", st.LibraryName, delimiter)
	}

	fmt.Fprintf(&msg, "BX%s%c", st.SupportedMessages.Marshal(), delimiter)

	if st.TerminalLocation != "" {
		fmt.Fprintf(&msg, "AN%s%c", st.TerminalLocation, delimiter)
	}

	if st.ScreenMessage != "" {
		fmt.Fprintf(&msg, "AF%s%c", st.ScreenMessage, delimiter)
	}

	if st.PrintLine != "" {
		fmt.Fprintf(&msg, "AG%s%c", st.PrintLine, delimiter)
	}

	if seqNum < 0 {
		seqNum = 0
	}

	return fmt.Sprintf("%s%c", utils.AppendChecksum(fmt.Sprintf("%sAY%dAZ", msg.String(), seqNum)), terminator)
}

func (st *ACSStatus) Unmarshal(line string, delimiter, terminator rune) (seqNum int, err error) {
	runes := []rune(line)

	if len(runes) < 42 {
		return 0, ErrInvalidResponse98
	}

	if string(runes[0:2]) != MsgIDACSStatus {
		return 0, ErrInvalidResponse98
	}

	codes := utils.ExtractFields(string(runes[36:]), delimiter, map[string]string{"AY": "", "AO": "", "AM": "", "BX": "", "AN": "", "AF": "", "AG": ""})
	seqNumString := codes["AY"]
	if seqNumString == "" {
		seqNum = 0
	} else {
		seqNum, err = strconv.Atoi(seqNumString)
		if err != nil {
			seqNum = 0
		}
	}

	st.OnlineStatus = utils.ParseBool(runes[2])
	st.CheckinOK = utils.ParseBool(runes[3])
	st.CheckoutOK = utils.ParseBool(runes[4])
	st.RenewalPolicy = utils.ParseBool(runes[5])
	st.StatusUpdateOK = utils.ParseBool(runes[6])
	st.OfflineOK = utils.ParseBool(runes[7])

	st.TimeoutPeriod, err = strconv.Atoi(string(runes[8:11]))
	if err != nil {
		return 0, err
	}

	st.RetriesAllowed, err = strconv.Atoi(string(runes[11:14]))
	if err != nil {
		return 0, err
	}

	st.DateTimeSync, err = time.Parse(utils.SIPDateFormat, string(runes[14:32]))
	if err != nil {
		return 0, err
	}

	st.ProtocolVersion = string(runes[32:36])

	st.InstitutionID = codes["AO"]

	if codes["AM"] != "" {
		st.LibraryName = codes["AM"]
	}

	st.SupportedMessages.Unmarshal(codes["BX"])

	if codes["AN"] != "" {
		st.TerminalLocation = codes["AN"]
	}

	if codes["AF"] != "" {
		st.ScreenMessage = codes["AF"]
	}

	if codes["AG"] != "" {
		st.PrintLine = codes["AG"]
	}

	err = st.Validate()
	if err != nil {
		return 0, err
	}

	return seqNum, nil
}

func (st *ACSStatus) Validate() error {
	err := Validate.Struct(st)
	if err != nil {
		return fmt.Errorf("invalid SIP %s response did not pass validation: %v", MsgIDACSStatus, err.(validator.ValidationErrors))
	}
	return nil
}

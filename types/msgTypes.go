package types

type MsgType int

const (
	_ MsgType = iota
	ReqBlockPatron
	_
	_
	_
	_
	_
	_
	_
	ReqCheckin
	RespCheckin
	ReqCheckout
	RespCheckout
	_
	_
	ReqHold
	RespHold
	ReqItemInfo
	RespItemInfo
	ReqItemStatusUpdate
	RespItemStatusUpdate
	_
	_
	ReqPatronStatus
	RespPatronStatus
	ReqPatronEnable
	RespPatronEnable
	_
	_
	ReqRenew
	RespRenew
	_
	_
	_
	_
	ReqEndPatronSession
	RespEndSession
	ReqFeePaid
	RespFeePaid
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	ReqPatronInfo
	RespPatronInfo
	ReqRenewAll
	RespRenewAll
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	ReqSCLogin
	RespSCLogin
	_
	RespSCResend
	ReqACSResend
	RespACSStatus
	ReqSCStatus
)

var msgIDs = [...]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "90", "91", "92", "93", "94", "95", "96", "97", "98", "99"}

var msgTypes = [...]string{
	"",
	"Block Patron Request",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"Checkin Request",
	"Checkin Response",
	"Checkout Request",
	"Checkout Response",
	"",
	"",
	"Hold Request",
	"Hold Response",
	"Item Info Request",
	"Item Info Response",
	"Item Status Update Request",
	"Item Status Update Response",
	"",
	"",
	"Patron Status Request",
	"Patron Status Response",
	"Patron Enable Request",
	"Patron Enable Response",
	"",
	"",
	"Renew Request",
	"Renew Response",
	"",
	"",
	"",
	"",
	"End Patron Session Request",
	"End Session Response",
	"Fee Paid Request",
	"Fee Paid Response",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"Patron Info Request",
	"Patron Info Response",
	"Renew All Request",
	"Renew All Response",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"SC Login Request",
	"SC Login Response",
	"",
	"SC Resend Response",
	"ACS Resend Request",
	"ACS Status Response",
	"SC Status Request",
}

func (m MsgType) ID() string {
	return msgIDs[m]
}

func (m MsgType) String() string {
	return msgTypes[m]
}

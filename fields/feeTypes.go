package fields

type FeeType int

const (
	_ FeeType = iota
	OtherUnknownFee
	AdministrativeFee
	DamageFee
	OverdueFee
	ProcessingFee
	RentalFee
	ReplacementFee
	ComputerAccessChargeFee
	HoldFee
)

var feeIDs = [...]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "90", "91", "92", "93", "94", "95", "96", "97", "98", "99"}

var feeTypes = [...]string{
	"Unknown Fee Type",
	"Other/Unknown",
	"Administrative",
	"Damage",
	"Overdue",
	"Processing",
	"Rental",
	"Replacement",
	"Computer Access Charge",
	"Hold Fee",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
	"Unknown Fee Type",
}

func (f FeeType) ID() string {
	if int(f) >= len(feeIDs) || f < 0 {
		return feeIDs[0]
	}
	return feeIDs[f]
}

func (f FeeType) String() string {
	if int(f) >= len(feeTypes) || f < 0 {
		return "Unknown Fee Type"
	}
	return feeTypes[f]
}

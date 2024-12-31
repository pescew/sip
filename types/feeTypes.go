package types

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

var feeIDs = [...]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09"}

var feeTypes = [...]string{
	"",
	"Other/Unknown",
	"Administrative",
	"Damage",
	"Overdue",
	"Processing",
	"Rental",
	"Replacement",
	"Computer Access Charge",
	"Hold Fee",
}

func (f FeeType) ID() string {
	if int(f) >= len(feeIDs) || int(f) < 1 {
		return feeIDs[1]
	}
	return feeIDs[f]
}

func (f FeeType) String() string {
	if int(f) >= len(feeTypes) || int(f) < 1 {
		return feeTypes[1]
	}
	return feeTypes[f]
}

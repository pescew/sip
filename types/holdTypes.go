package types

type HoldType int

const (
	_ HoldType = iota
	HoldTypeOther
	HoldTypeAny
	HoldTypeSpecific
	HoldTypeAnyAtLocation
)

var holdIDs = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var holdTypes = [...]string{
	"Unknown hold type",
	"Other",
	"Any copy of a title",
	"A specific copy of a title",
	"Any copy at a single branch or sublocation",
	"Unknown hold type",
	"Unknown hold type",
	"Unknown hold type",
	"Unknown hold type",
	"Unknown hold type",
}

func (h HoldType) ID() string {
	if int(h) >= len(holdIDs) || h < 0 {
		return holdIDs[0]
	}
	return holdIDs[h]
}

func (h HoldType) String() string {
	if int(h) >= len(holdTypes) || h < 0 {
		return holdTypes[0]
	}
	return holdTypes[h]
}

package types

type HoldType int

const (
	_ HoldType = iota
	HoldTypeOther
	HoldTypeAny
	HoldTypeSpecific
	HoldTypeAnyAtLocation
)

var HoldIDs = [...]string{"0", "1", "2", "3", "4"}

var HoldTypes = [...]string{
	"",
	"Other",
	"Any copy of a title",
	"A specific copy of a title",
	"Any copy at a single branch or sublocation",
}

func (h HoldType) ID() string {
	if int(h) >= len(HoldTypes) || int(h) < 1 {
		return HoldIDs[1]
	}
	return HoldIDs[h]
}

func (h HoldType) String() string {
	if int(h) >= len(HoldTypes) || int(h) < 1 {
		return HoldTypes[1]
	}
	return HoldTypes[h]
}

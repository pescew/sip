package types

type HoldType int

const (
	_ HoldType = iota
	HoldTypeOther
	HoldTypeAny
	HoldTypeSpecific
	HoldTypeAnyAtLocation
)

var holdIDs = [...]string{"0", "1", "2", "3", "4"}

var holdTypes = [...]string{
	"",
	"Other",
	"Any copy of a title",
	"A specific copy of a title",
	"Any copy at a single branch or sublocation",
}

func (h HoldType) ID() string {
	if int(h) >= len(holdIDs) || int(h) < 1 {
		return holdIDs[1]
	}
	return holdIDs[h]
}

func (h HoldType) String() string {
	if int(h) >= len(holdTypes) || int(h) < 1 {
		return holdTypes[1]
	}
	return holdTypes[h]
}

package types

type SCStatusCode int

const (
	SCStatusOK SCStatusCode = iota
	SCStatusOutOfPaper
	SCStatusShuttingDown
)

var scStatusIDs = [...]string{"0", "1", "2"}

var scStatusCodes = [...]string{
	"Selfcheck unit is OK",
	"Selfcheck pritner is out of paper",
	"Selfcheck is about to shut down",
}

func (s SCStatusCode) ID() string {
	if int(s) >= len(scStatusIDs) || int(s) < 1 {
		return scStatusIDs[1]
	}
	return scStatusIDs[s]
}

func (s SCStatusCode) String() string {
	if int(s) >= len(scStatusCodes) || int(s) < 1 {
		return scStatusCodes[1]
	}
	return scStatusCodes[s]
}

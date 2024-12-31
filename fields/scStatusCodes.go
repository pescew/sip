package fields

type SCStatusCode int

const (
	SCStatusOK SCStatusCode = iota
	SCStatusOutOfPaper
	SCStatusShuttingDown
)

var scStatusIDs = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var scStatusCodes = [...]string{
	"Selfcheck unit is OK",
	"Selfcheck pritner is out of paper",
	"Selfcheck is about to shut down",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
	"Unknown selfcheck status code",
}

func (s SCStatusCode) ID() string {
	if int(s) >= len(scStatusIDs) || s < 0 {
		return scStatusIDs[0]
	}
	return scStatusIDs[s]
}

func (s SCStatusCode) String() string {
	if int(s) >= len(scStatusCodes) || s < 0 {
		return "Unknown selfcheck status code"
	}
	return scStatusCodes[s]
}

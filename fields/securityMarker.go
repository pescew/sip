package fields

type SecurityMarker int

const (
	SecurityMarkerOther SecurityMarker = iota
	SecurityMarkerNone
	SecurityMarker3MTattleTapeSecurityStrip
	SecurityMarker3MWhisperTape
)

var securityMarkerIDs = [...]string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "90", "91", "92", "93", "94", "95", "96", "97", "98", "99"}

var securityMarkerTypes = [...]string{
	"Other",
	"None",
	"3M Tattle-Tape Security Strip",
	"3M Whisper Tape",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
	"Unknown Security Marker",
}

func (s SecurityMarker) ID() string {
	if int(s) >= len(securityMarkerIDs) || s < 0 {
		return securityMarkerIDs[0]
	}
	return securityMarkerIDs[s]
}

func (s SecurityMarker) String() string {
	if int(s) >= len(securityMarkerTypes) || s < 0 {
		return "Unknown Security Marker"
	}
	return securityMarkerTypes[s]
}

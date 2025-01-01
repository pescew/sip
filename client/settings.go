package client

type Config struct {
	DebugMode                bool
	LibraryID                string
	InstitutionID            string
	TerminalUsername         string
	TerminalPassword         string
	TerminatorCharacter      rune
	DelimiterCharacter       rune
	ConnectionTimeoutSeconds int
	ErrorDetection           bool
}

func DefaultConfig() Config {
	return Config{
		DebugMode:                false,
		LibraryID:                "lib",
		InstitutionID:            "inst",
		TerminalUsername:         "",
		TerminalPassword:         "",
		TerminatorCharacter:      '\r',
		DelimiterCharacter:       '|',
		ConnectionTimeoutSeconds: 5,
		ErrorDetection:           true,
	}
}

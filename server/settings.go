package server

type Config struct {
	Host                string
	Port                int
	DebugMode           bool
	LibraryID           string
	InstitutionID       string
	TerminalUsername    string
	TerminalPassword    string
	TerminatorCharacter rune
	DelimiterCharacter  rune
	ConnectionTimeout   int
	ErrorDetection      bool
}

func DefaultConfig() Config {
	return Config{
		Host:                "127.0.0.1",
		Port:                9000,
		DebugMode:           false,
		LibraryID:           "lib",
		InstitutionID:       "inst",
		TerminalUsername:    "",
		TerminalPassword:    "",
		TerminatorCharacter: '\r',
		DelimiterCharacter:  '|',
		ConnectionTimeout:   5,
		ErrorDetection:      true,
	}
}

type Settings struct {
	host                string
	port                int
	debugMode           bool
	libraryID           string
	institutionID       string
	terminalUsername    string
	terminalPassword    string
	terminatorCharacter rune
	delimiterCharacter  rune
	connectionTimeout   int
	errorDetection      bool
}

func (s *Settings) Host() string {
	return s.host
}

func (s *Settings) Port() int {
	return s.port
}

func (s *Settings) DebugMode() bool {
	return s.debugMode
}

func (s *Settings) LibraryID() string {
	return s.libraryID
}

func (s *Settings) InstitutionID() string {
	return s.institutionID
}

func (s *Settings) TerminalUsername() string {
	return s.terminalUsername
}

func (s *Settings) TerminalPassword() string {
	return s.terminalPassword
}

func (s *Settings) TerminatorCharacter() rune {
	return s.terminatorCharacter
}

func (s *Settings) DelimiterCharacter() rune {
	return s.delimiterCharacter
}

func (s *Settings) ConnectionTimeout() int {
	return s.connectionTimeout
}

func (s *Settings) ErrorDetection() bool {
	return s.errorDetection
}

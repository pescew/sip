package types

type Algorithm int

const (
	AlgorithmPlaintext Algorithm = iota
	_
	_
	_
	_
	_
	_
	_
	_
	_
)

var algorithmIDs = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var algorithmDescriptions = [...]string{
	"Unencrypted plaintext",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
	"Unknown algorithm",
}

func (s Algorithm) ID() string {
	if int(s) >= len(algorithmIDs) || s < 0 {
		return algorithmIDs[0]
	}
	return algorithmIDs[s]
}

func (s Algorithm) String() string {
	if int(s) >= len(algorithmDescriptions) || s < 0 {
		return algorithmDescriptions[0]
	}
	return algorithmDescriptions[s]
}

package fields

type HoldMode string

const (
	HoldModeAdd    HoldMode = "+"
	HoldModeDelete HoldMode = "-"
	HoldModeChange HoldMode = "*"
)

var holdModes = map[string]string{
	"+": "Add patron to the hold queue for the item",
	"-": "Delete patron from the hold queue for the item",
	"*": "Change the hold to match the message parameters",
}

func (h HoldMode) ID() string {
	return string(h)
}

func (h HoldMode) String() string {
	switch h {
	case "+":
		return holdModes["+"]
	case "-":
		return holdModes["-"]
	case "*":
		return holdModes["*"]
	default:
		return "Unknown hold mode"
	}
}

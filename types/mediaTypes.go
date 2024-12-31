package types

type MediaType string

const (
	MediaTypeOther             MediaType = "000"
	MediaTypeBook              MediaType = "001"
	MediaTypeMagazine          MediaType = "002"
	MediaTypeBoundJournal      MediaType = "003"
	MediaTypeAudioTape         MediaType = "004"
	MediaTypeVideoTape         MediaType = "005"
	MediaTypeCDROM             MediaType = "006"
	MediaTypeDiskette          MediaType = "007"
	MediaTypeBookWithDiskette  MediaType = "008"
	MediaTypeBookWithCD        MediaType = "009"
	MediaTypeBookWithAudioTape MediaType = "010"
)

var mediaTypes = map[string]string{
	"000": "Other",
	"001": "Book",
	"002": "Magazine",
	"003": "Bound Journal",
	"004": "Audio Tape",
	"005": "Video Tape",
	"006": "CD/CDROM",
	"007": "Diskette",
	"008": "Book with Diskette",
	"009": "Book with CD",
	"010": "Book with Audio Tape",
}

func (m MediaType) ID() string {
	return string(m)
}

func (m MediaType) String() string {
	val, exists := mediaTypes[string(m)]
	if exists {
		return val
	}
	return "Unknown Media Type"
}

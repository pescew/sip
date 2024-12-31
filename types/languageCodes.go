package types

import "fmt"

type LanguageCode int

const (
	LanguageUnknown LanguageCode = iota
	LanguageEnglish
	LanguageFrench
	LanguageGerman
	LanguageItalian
	LanguageDutch
	LanguageSwedish
	LanguageFinnish
	LanguageSpanish
	LanguageDanish
	LanguagePortuguese
	LanguageCanadianFrench
	LanguageNorwegian
	LanguageHebrew
	LanguageJapanese
	LanguageRussian
	LanguageArabic
	LanguagePolish
	LanguageGreek
	LanguageChinese
	LanguageKorean
	LanguageNorthAmericanSpanish
	LanguageTamil
	LanguageMalay
	LanguageUnitedKingdom
	LanguageIcelandic
	LanguageBelgian
	LanguageTaiwanese
)

var languageIDs = [...]string{"000", "001", "002", "003", "004", "005", "006", "007", "008", "009", "010", "011", "012", "013", "014", "015", "016", "017", "018", "019", "020", "021", "022", "023", "024", "025", "026", "027"}

var languageCodes = [...]string{
	"Unknown",
	"English",
	"French",
	"German",
	"Italian",
	"Dutch",
	"Swedish",
	"Finnish",
	"Spanish",
	"Danish",
	"Portuguese",
	"Canadian-French",
	"Norwegian",
	"Hebrew",
	"Japanese",
	"Russian",
	"Arabic",
	"Polish",
	"Greek",
	"Chinese",
	"Korean",
	"North American Spanish",
	"Tamil",
	"Malay",
	"United Kingdom",
	"Icelandic",
	"Belgian",
	"Taiwanese",
}

func (l LanguageCode) ID() string {
	if l < 0 || l > 999 {
		return languageIDs[0]
	} else if int(l) >= len(languageIDs) {
		return fmt.Sprintf("%03d", l)
	}
	return languageIDs[l]
}

func (l LanguageCode) String() string {
	if int(l) >= len(languageCodes) || l < 0 {
		return languageCodes[0]
	}
	return languageCodes[l]
}

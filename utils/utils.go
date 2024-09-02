package utils

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

const (
	SIPDateFormat = "20060102    150405"
	// SIPMaxFieldsPerRequest = 30
	// SIPMaxItemsPerRequest  = 100
)

var REPLACER *strings.Replacer

func EscapeSIP(text string) string {
	return REPLACER.Replace(text)
}

func ConfigureEscapeCharacters(chars ...rune) {
	replace := []string{}
	for _, char := range chars {
		replace = slices.Concat(replace, []string{string(char), ""})
	}
	REPLACER = strings.NewReplacer(replace...)
}

func YorN(field bool) string {
	if field {
		return "Y"
	}
	return "N"
}

func YorBlank(field bool) string {
	if field {
		return "Y"
	}
	return " "
}

func ZeroOrOne(field bool) string {
	if field {
		return "1"
	}
	return "0"
}

func ParseBool(char rune) bool {
	switch string(char) {
	case "Y", "y", "1":
		return true
	}
	return false
}

func ComputeChecksum(msg string) string {
	check := 0
	for _, character := range msg {
		check += int(character)
	}
	check += int('\x00') //null terminate string
	check = (check ^ 0xFFFF) + 1
	checksum := fmt.Sprintf("%4.4X", check)
	return checksum
}

func AppendChecksum(msg string) string {
	return msg + ComputeChecksum(msg)
}

func GenerateLineScanner(terminator rune) func([]byte, bool) (int, []byte, error) {
	terminatorBytes := []byte(string(terminator))
	terminatorLength := len(terminatorBytes)

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, terminatorBytes); i >= 0 {
			return i + terminatorLength, data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}

func GenerateSIPValidatorFunc(badChars string) func(validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		fieldString := fl.Field().String()
		if utf8.RuneCountInString(fieldString) > 255 {
			return false
		}
		return !strings.ContainsAny(fieldString, badChars)
	}
}

func ExtractFields(line string, delimiter rune, fields map[string]string) map[string]string {
	var segment string
	delim := string(delimiter)
	count := 0
	found := true
	for found {
		segment, line, found = strings.Cut(line, delim)
		seg := []rune(segment)
		if len(seg) > 2 {
			field := string(seg[0:2])
			_, exists := fields[field]
			if exists {
				if field == "AY" {
					fields[field] = string(seg[2:3])
				} else {
					fields[field] = string(seg[2:])
				}
			}
		}
		// if count > SIPMaxFieldsPerRequest {
		// 	break
		// }
		count++
	}

	return fields
}

func ExtractMultiFields(line string, delimiter rune, fields map[string][]string) map[string][]string {
	var segment string
	delim := string(delimiter)
	count := 0
	found := true
	for found {
		segment, line, found = strings.Cut(line, delim)
		seg := []rune(segment)
		if len(seg) > 2 {
			field := string(seg[0:2])
			_, exists := fields[field]
			if exists {
				fields[field] = append(fields[field], string(seg[2:]))
			}
		}
		// if count > SIPMaxItemsPerRequest {
		// 	break
		// }
		count++
	}

	return fields
}

func IncrementSeqNum(seqNum int) int {
	if seqNum == 9 || seqNum < 0 {
		return 0
	}
	return seqNum + 1
}

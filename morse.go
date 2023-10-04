package morse

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

// Phrase represets message in Morse code.
type Phrase string

// Translator interface translates the runes in the string into
// a set of characters representing this rune in Morse code.
// It is necessary for the formation of phrases in Morse code.
// Spaces must be translated into "/",
// The r rune must be translated into a set dots ('.') and dashes('-')
type Translator interface {
	// translates the rune in the string into
	// a set of characters representing this rune in Morse code.
	// If the rune cannot be translated, it returns an empty string.
	Translate(r rune) string
}

type jsonTranslator struct {
	morseMap map[rune]string
}

// Func JSONTranslator returns a Translator that takes a table
// of characters from a json file along the path "path".
func JSONTranslator(path string) (*jsonTranslator, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, errors.New("JSONTranslator initialization error: " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New("JSONTranslator initialization error: " + err.Error())
	}

	var morseMap map[rune]string
	err = json.Unmarshal(byteValue, &morseMap)
	if err != nil {
		err = errors.New("JSONTranslator initialization error: " + err.Error())
	}

	return &jsonTranslator{morseMap}, err
}

func (t jsonTranslator) Translate(r rune) string {
	u := unicode.ToUpper(r)
	return t.morseMap[u]
}

// Converts a string using a translator into a phrase in Morse code.
func Parse(s string, t Translator) Phrase {
	var b bytes.Buffer
	for _, r := range s {
		code := t.Translate(r)
		b.WriteString(code)
		b.WriteByte(' ')
	}

	str := b.String()
	str = strings.Trim(str, " ")

	return Phrase(str)
}

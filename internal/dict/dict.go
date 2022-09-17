package dict

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

// language
var (
	Russian  = "ru"
	Ukranian = "uk"
	English  = "en"
	Spanish  = "es"
)

// lexis kinds, idk why
var (
	Noun         = "noun"
	Pronoun      = "pronoun"
	Verb         = "verb"
	Adjective    = "adjective"
	Adverb       = "adverb"
	Preposition  = "preposition"
	Conjunction  = "conjunction"
	Interjection = "interjection"
)

// representation of word in all languages
type Word struct {
	English  Language `json:"en"`
	Russian  Language `json:"ru,omitempty"`
	Ukranian Language `json:"uk,omitempty"`
	Spanish  Language `json:"es,omitempty"`

	Type string `json:"type"`
}

type Language struct {
	Transcription string `json:"transcription,omitempty"`

	Word string `json:"word,omitempty"`

	Definition string `json:"definition,omitempty"`
}

type Dictionary []Word

func New850() (Dictionary, error) {
	b, err := os.ReadFile("850.json")
	if err != nil {
		return nil, err
	}
	var d Dictionary
	err = json.Unmarshal(b, &d)
	return d, err
}

// Export dictionary to json file. No need to add extension.
func (d Dictionary) ToFile(name string) {
	f, err := os.Create(name + ".json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
}

func (d Dictionary) WordByLanguage(lang, word string) Word {
	for _, w := range d {
		switch lang {
		case English:
			if w.English.Word == word {
				return w
			}
		case Russian:
			if w.Russian.Word == word {
				return w
			}
		case Spanish:
			if w.Spanish.Word == word {
				return w
			}
		}
	}
	return Word{}
}

func (d Dictionary) FileByLanguage(lang string) error {
	f, err := os.Create(lang + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()
	var words string
	for _, w := range d {
		switch lang {
		case English:
			words += w.English.Word + "\n"
		case Russian:
			words += w.Russian.Word + "\n"
		case Spanish:
			words += w.Spanish.Word + "\n"
		case Ukranian:
			words += w.Ukranian.Word + "\n"
		}
	}

	_, err = f.Write([]byte(words))
	return err
}

func (w *Word) Fill(lang, word string) {
	switch lang {
	case English:
		w.English.Word = word
	case Russian:
		w.Russian.Word = word
	case Spanish:
		w.Spanish.Word = word
	case Ukranian:
		w.Ukranian.Word = word
	}
}

// NOTE: this is not a real code, just a snippet
func GetSpelling(lang, word string) string {
	url := "https://api2.unalengua.com/ipav3"

	//{
	// 	"text": word,
	// 	"lang": lang+"-US",
	// 	"mode": true
	//}

	req, err := http.NewRequest("POST", url, strings.NewReader(`{"text":"`+word+`","lang":"`+lang+`-US","mode":true}`))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// {"detected":"UNKNOWN","ipa":"kwˌaðɾikˈopːteɾo","lang":"es-US","spelling":"cuadricóptero"}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	return m["ipa"].(string)
}

// File should be with 850 words separated by new line
func (d Dictionary) FillFromFile(lang, filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file := string(b)
	words := strings.Split(file, "\n")

	if len(words) == 851 {
		words = words[:850]
	}

	for i := range d {
		d[i].Fill(Ukranian, words[i])
	}

	return nil
}
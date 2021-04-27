package api

type Language struct {
	Code           string
	Names          map[string]string
	Id             int
	Transliterated bool
}

type LanguageResponse struct {
	Languages []Language
}

type Example struct {
	First_phrase_bold_offsets [][]int
	First_phrase              string
	Second_phrase             string
	Equals                    bool
	Sirst_phrase_bold_offsets [][]int
}

type Pronunciation struct {
	Dialect string
	Ipa     string
	Mp3url  string
}

type Word struct {
	Context         string
	Pronunciation   []Pronunciation
	Language        int
	Id              int
	Text            string
	Transliteration string
	Gender          string
}

type Synonym struct {
	Language int
	Id       int
	Text     string
}

type Translation struct {
	Contexts []string
	Word     Word
}

type Group struct {
	Name         string
	Translations []Translation
}

type Definition struct {
	Definitions []string
	Group       string
}

type Inflection struct {
	Type string
	Word Word
}

type Suggestion struct {
	Word      Word
	Relevance float32
}

type SearchResponse struct {
	Alternative_spellings []Word
	Examples              []Example
	Synonyms              []Synonym
	Suggestions           []Suggestion
	Groups                []Group
	Relations             []Inflection
	Word                  Word
	Definitions           []Definition
	Inflections           []Inflection
	Abbreviations         []Word
	Status                string
}

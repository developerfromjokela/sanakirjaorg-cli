package main

import (
	"encoding/base64"
	json2 "encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"regexp"
	"sanakirjaorg-cli/api"
	"sanakirjaorg-cli/log"
	"sanakirjaorg-cli/utils"
	"strings"
)

var interactive bool
var json bool
var version bool
var prettify bool
var pronun bool
var altSpellings bool
var examples bool
var definitions bool
var synonyms bool
var verbose int
var command string
var source string
var target string
var text string

var languages *api.LanguageResponse

var defaultLang string

func main() {
	flag.BoolVar(&interactive, "i", false, "Interactive mode")
	flag.BoolVar(&version, "v", false, "Version")
	flag.BoolVar(&examples, "exp", false, "Show examples")
	flag.BoolVar(&json, "json", false, "Print all data as JSON")
	flag.BoolVar(&prettify, "prettify", false, "Prettify JSON data")
	flag.BoolVar(&pronun, "pro", false, "Show pronunciation")
	flag.BoolVar(&altSpellings, "alt", false, "Show alternative spellings")
	flag.BoolVar(&synonyms, "syn", false, "Show synonyms")
	flag.BoolVar(&definitions, "def", false, "Show definitions")
	flag.IntVar(&verbose, "vv", log.WARN, "Verbose mode number")
	flag.StringVar(&command, "c", "", "Command, such as: langlist")
	flag.StringVar(&defaultLang, "l", "en", "Language, available options: fi, sv, en, fr")
	flag.StringVar(&source, "src", "", "Source language (code), full list via command langlist")
	flag.StringVar(&target, "target", "", "Target language (code), full list via command langlist")
	flag.StringVar(&text, "text", "", "Your word to define")
	flag.Parse()
	var logLevel = log.WARN
	verboseVal := verbose
	logLevel = verboseVal
	log.LogLevel = logLevel
	if !validateLang() {
		log.Log(log.WARN, "Invalid language option: "+defaultLang)
		flag.Usage()
		return
	}
	if version {
		var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(utils.AsciiART)))
		decodePos, err := base64.StdEncoding.Decode(decoded, []byte(utils.AsciiART))
		if err != nil {
			log.Log(log.ERROR, err.Error())
			return
		}
		fmt.Printf(utils.Red, string(decoded[:decodePos]))
		fmt.Printf(utils.Yellow, "Sanakirja.org CLI - v"+utils.Version)
		fmt.Println()
		fmt.Printf(utils.Yellow, "Github: https://github.com/developerfromjokela/sanakirjaorg-cli")
		fmt.Println()
		return
	}
	if !interactive {
		nonInteractive()
	} else {
		interactiveMode()
	}
}

func validateLang() bool {
	return defaultLang == "fi" || defaultLang == "sv" || defaultLang == "en" || defaultLang == "fr"
}

func nonInteractive() {
	if isFlagPassed("c") {
		if command == "langlist" {
			fetchLanguages()
			if !json {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Code", "Name", "Transliterated"})
				for _, lang := range languages.Languages {
					if lang.Code != "null" {
						table.Append([]string{lang.Code, lang.Names[defaultLang], transliteration(lang.Transliterated)})
					}
				}
				table.Render()
			} else {
				jsonExport(languages)
			}
			return
		}
		flag.Usage()
	} else if isFlagPassed("src") && isFlagPassed("target") && isFlagPassed("text") {
		fetchLanguages()
		var sourceLang = findLanguage(source)
		var targetLang = findLanguage(target)
		nonInteractiveSearch(sourceLang, targetLang)
	} else {
		log.Log(log.DEBUG, len(flag.Args()))
		log.Log(log.DEBUG, flag.Args())
		if len(flag.Args()) == 1 {
			fmt.Println("Using defaults: en->fi")
			fetchLanguages()
			var sourceLang = findLanguage("en")
			var targetLang = findLanguage("fi")
			text = flag.Arg(0)
			nonInteractiveSearch(sourceLang, targetLang)
			return
		} else if len(flag.Args()) == 2 {
			fmt.Println("Using default target (fi)")
			fetchLanguages()
			text = flag.Arg(1)
			var sourceLang = findLanguage(flag.Arg(0))
			var targetLang = findLanguage("fi")
			nonInteractiveSearch(sourceLang, targetLang)
			return
		} else if len(flag.Args()) == 3 {
			fetchLanguages()
			text = flag.Arg(2)
			var sourceLang = findLanguage(flag.Arg(0))
			var targetLang = findLanguage(flag.Arg(1))
			nonInteractiveSearch(sourceLang, targetLang)
			return
		}
		// HELP
		flag.Usage()
	}
}

func nonInteractiveSearch(sourceLang *api.Language, targetLang *api.Language) {
	if sourceLang == nil {
		log.Log(log.WARN, "Source language is invalid")
		flag.Usage()
		return
	}
	if targetLang == nil {
		log.Log(log.WARN, "Source language is invalid")
		flag.Usage()
		return
	}
	searchResp, err := api.Define(sourceLang.Id, targetLang.Id, text, defaultLang)
	if err != nil {
		log.Log(log.ERROR, err)
		return
	}
	if searchResp.Status == "ok" {
		if !json {
			// Basic mode
			for _, group := range searchResp.Groups {
				fmt.Println("*** " + group.Name + " ***")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Word", "Pronountiation"})
				for _, translation := range group.Translations {
					var pronounciation []string
					for _, pronoun := range translation.Word.Pronunciation {
						pronounciation = append(pronounciation, pronoun.Ipa)
					}
					table.Append([]string{translation.Word.Text, strings.Join(pronounciation, ",")})
				}
				table.Render()
			}
			if definitions {
				fmt.Println()
				fmt.Println("--> Definitions")
				for _, definition := range searchResp.Definitions {
					fmt.Println("** ", definition.Group, " **")
					for pos, item := range definition.Definitions {
						var correctedDefinition = item
						correctedDefinition = strings.ReplaceAll(correctedDefinition, "<em>", "")
						correctedDefinition = strings.ReplaceAll(correctedDefinition, "</em>", "")
						var re = regexp.MustCompile(`\[\[(.*?)\|(.*?)\]\]`)
						correctedDefinition = re.ReplaceAllString(correctedDefinition, `$2`)
						fmt.Println("[", pos+1, "] ", correctedDefinition)
					}
				}
				fmt.Println()
			}
			if synonyms {
				fmt.Println()
				fmt.Println("--> Synonyms")
				for _, example := range searchResp.Synonyms {
					fmt.Println(example.Text)
				}
				fmt.Println()
			}
			if altSpellings {
				fmt.Println()
				fmt.Println("--> Alternative spellings")
				for _, example := range searchResp.Alternative_spellings {
					fmt.Print(example.Text)
					if len(example.Context) > 0 {
						fmt.Print(" ", "(", example.Context, ")")
					}
					fmt.Println("")
				}
				fmt.Println()
			}
			if examples {
				fmt.Println()
				fmt.Println("--> Examples")
				for _, example := range searchResp.Examples {
					fmt.Println("[*] " + example.First_phrase)
					if len(example.Second_phrase) > 0 {
						fmt.Println("[*] " + example.Second_phrase)
					}
				}
				fmt.Println()
			}
			if examples {
				fmt.Println()
				fmt.Println("--> Examples")
				for _, example := range searchResp.Examples {
					fmt.Println("[*] " + example.First_phrase)
					if len(example.Second_phrase) > 0 {
						fmt.Println("[*] " + example.Second_phrase)
					}
				}
				fmt.Println()
			}
			if pronun {
				fmt.Println()
				fmt.Println("--> Pronunciations")
				for _, pronun := range searchResp.Word.Pronunciation {
					var dialect = ""
					if pronun.Dialect != "*" {
						dialect = pronun.Dialect + " "
					}
					fmt.Println("[*] ", dialect, pronun.Ipa)
				}
				fmt.Println()
			}
		} else {
			// Export whole thing as JSON
			jsonExport(searchResp)
		}
	} else if searchResp.Status == "error" {
		log.Log(log.ERROR, "Sanakirja.org server error")
	} else if searchResp.Status == "word_not_found" || searchResp.Status == "no_translations_found" {
		fmt.Print("Word you're searching was not found.")
		if searchResp.Suggestions != nil && len(searchResp.Suggestions) > 0 {
			fmt.Print(" Here are some suggestions:")
		}
		fmt.Println("")
		for _, suggestion := range searchResp.Suggestions {
			fmt.Println(suggestion.Word.Text)
		}
	} else {
		log.Log(log.LOG, "Status not found! FATAL: "+searchResp.Status)
		log.Log(log.ERROR, "Parsing error. Please enable verbose to see exact cause")
	}
}

func jsonExport(item interface{}) {
	// Export whole thing as JSON
	if !prettify {
		b, err := json2.Marshal(&item)
		if err != nil {
			log.Log(log.ERROR, err.Error())
			return
		}
		fmt.Println(string(b))
	} else {
		b, err := json2.MarshalIndent(&item, "", "\t")
		if err != nil {
			log.Log(log.ERROR, err.Error())
			return
		}
		fmt.Println(string(b))
	}
}
func transliteration(available bool) string {
	if available {
		return "Yes"
	}
	return "No"
}

func findLanguage(text string) *api.Language {
	var language *api.Language = nil
	for _, lang := range languages.Languages {
		if strings.ToLower(text) == lang.Code {
			language = &lang
			break
		}
	}
	return language
}

func interactiveMode() {
	// TODO
	log.Log(log.WARN, "Interactive mode not available yet!")
}

func fetchLanguages() {
	langResp, err := api.Languages()
	if err != nil {
		log.Log(log.ERROR, err)
		return
	}
	languages = langResp
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

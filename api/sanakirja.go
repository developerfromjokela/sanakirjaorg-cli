package api

import (
	url2 "net/url"
	"sanakirjaorg-cli/log"
	"sanakirjaorg-cli/net"
	"strconv"
)
import "sanakirjaorg-cli/cache"

var Verbose *bool

func Languages() (*LanguageResponse, error) {
	var langRes = &LanguageResponse{}
	err := cache.GetCache(languageList, langRes)
	if err != nil {
		log.Log(log.DEBUG, err.Error())
		if err.Error() != "Cache expired" {
			log.Log(log.LOG, "Cache not present, re-generating")
		}
	}
	err = networking.GetJson(languageList, langRes)
	err2 := cache.SaveCache(languageList, langRes)
	if err2 != nil {
		log.Log(log.WARN, "Caching failed! Check file permissions to your home directory")
	}
	return langRes, err
}

func Define(sourceLanguage int, targetLanguage int, text string, locale string) (*SearchResponse, error) {
	var searchRes = &SearchResponse{}
	url := search + "&text=" + url2.QueryEscape(text) + "&sourceLanguage=" + strconv.Itoa(sourceLanguage) + "&targetLanguage=" + strconv.Itoa(targetLanguage) + "&locale=" + locale
	log.Log(log.DEBUG, "URL: "+url)
	err := cache.GetCache(url, searchRes)
	if err != nil {
		log.Log(log.DEBUG, err.Error())
		if err.Error() != "Cache expired" {
			log.Log(log.LOG, "Cache not present, re-generating")
		}
	}
	err = networking.GetJson(url, searchRes)
	err2 := cache.SaveCache(url, searchRes)
	if err2 != nil {
		log.Log(log.WARN, "Caching failed! Check file permissions to your home directory")
	}
	return searchRes, err
}

func DefineWithWordId(sourceLanguage int, targetLanguage int, text int, locale string) (*SearchResponse, error) {
	var searchRes = &SearchResponse{}
	url := search + "&wordId=" + string(rune(text)) + "&sourceLanguage=" + string(rune(sourceLanguage)) + "&targetLanguage=" + string(rune(targetLanguage)) + "&locale=" + locale
	err := cache.GetCache(url, searchRes)
	if err != nil {
		log.Log(log.DEBUG, err.Error())
		if err.Error() != "Cache expired" {
			log.Log(log.LOG, "Cache not present, re-generating")
		}
	}
	err = networking.GetJson(url, searchRes)
	err2 := cache.SaveCache(url, searchRes)
	if err2 != nil {
		log.Log(log.WARN, "Caching failed! Check file permissions to your home directory")
	}
	return searchRes, err
}

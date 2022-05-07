package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Translate struct using for google translate api
type Translate struct {
	Sentences []struct {
		Backend int64  `json:"backend"`
		Orig    string `json:"orig"`
		Trans   string `json:"trans"`
	} `json:"sentences"`
	Spell struct{} `json:"spell"`
	Src   string   `json:"src"`
}

// tr(string, string, string) string
// Simple translate function with google http api
func tr(from, to, query string) string {
	newText := ""
	text := sprintf(gtranslate, from, to, url.QueryEscape(query))
	resp, err := http.Get(text)
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
		return query
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
		return query
	}
	var ret Translate
	if err = json.Unmarshal(body, &ret); err != nil {
		toLog(sprintf("!!! %s\n", err.Error()))
		return query
	}
	for _, v := range ret.Sentences {
		newText += v.Trans
	}
	if newText != "" {
		return newText
	}
	return query
}

// tr(string) function is wrapper around Message
// use default 'en' language in source code
func (o *Message) tr(text string) string {
	if o.LanguageCode == "en" {
		return text
	}
	return tr("en", o.LanguageCode, text)
}

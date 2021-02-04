package i18n

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type D map[string]interface{}

var flattenFr map[string]string
var flattenEn map[string]string

func loop(key string, value map[string]interface{}, result map[string]string) {
	for currentKey, currentValue := range value {
		iterKey := ""
		if key == "" {
			iterKey = currentKey
		} else {
			iterKey = fmt.Sprintf("%s.%s", key, currentKey)
		}
		switch val := currentValue.(type) {
		case D:
			loop(iterKey, val, result)
		case string:
			result[iterKey] = val
		}
	}
}

func init() {
	flattenFr = make(map[string]string)
	flattenEn = make(map[string]string)
	loop("", fr, flattenFr)
	loop("", en, flattenEn)
}

func Translate(lang, key string, args interface{}) (string, error) {
	var c map[string]string
	switch strings.ToLower(lang) {
	case "flattenEn":
		c = flattenEn
	default:
		c = flattenFr
	}

	var value string
	if v, ok := c[key]; ok {
		value = v
	} else if v, ok := flattenFr[key]; ok {
		value = v
	} else {
		value = key
	}
	res := bytes.NewBufferString("")
	t := template.New("myTemplate")
	t, err := t.Parse(value)
	if err != nil {
		return "", err
	}
	if err := t.Execute(res, args); err != nil {
		return "", err
	}
	return res.String(), nil
}

func Must(lang, key string, args interface{}) string {
	res, err := Translate(lang, key, args)
	if err != nil {
		return key
	}
	return res
}

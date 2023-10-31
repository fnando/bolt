package common

import (
	"os"
	"strings"
)

type color struct {
	Disabled bool
}

var defaultColors map[string]string = map[string]string{
	"text":   "30",
	"fail":   "31",
	"pass":   "32",
	"skip":   "33",
	"detail": "34",
}

func (c color) Apply(code string, text string) string {
	if c.Disabled {
		return text
	}

	return "\033[" + code + "m" + text + "\033[0m"
}

func (c color) Fail(text string) string {
	return c.Apply(c.Color("fail"), text)
}

func (c color) Text(text string) string {
	return c.Apply(c.Color("text"), text)
}

func (c color) Pass(text string) string {
	return c.Apply(c.Color("pass"), text)
}

func (c color) Skip(text string) string {
	return c.Apply(c.Color("skip"), text)
}

func (c color) Detail(text string) string {
	return c.Apply(c.Color("detail"), text)
}

func (c color) Color(name string) string {
	env := "BOLT_" + strings.ToUpper(name) + "_COLOR"
	val := os.Getenv(env)

	if val != "" {
		return val
	}

	val = defaultColors[strings.ToLower(name)]

	if val == "" {
		return "30"
	}

	return val
}

var Color color = color{}

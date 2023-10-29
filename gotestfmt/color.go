package gotestfmt

type color struct {
	TextColor   string
	FailColor   string
	PassColor   string
	SkipColor   string
	DetailColor string
	Disabled    bool
}

func (c color) apply(code string, text string) string {
	if c.Disabled {
		return text
	}

	return "\033[" + code + "m" + text + "\033[0m"
}

func (c color) Fail(text string) string {
	return c.apply(c.FailColor, text)
}

func (c color) Text(text string) string {
	return c.apply(c.TextColor, text)
}

func (c color) Pass(text string) string {
	return c.apply(c.PassColor, text)
}

func (c color) Skip(text string) string {
	return c.apply(c.SkipColor, text)
}

func (c color) Detail(text string) string {
	return c.apply(c.DetailColor, text)
}

var Color color = color{}

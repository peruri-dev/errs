package errs

import "fmt"

func NewCodex(title, detail, customCode string, httpCode int) *Codex {
	return &Codex{
		Title:      title,
		Detail:     detail,
		CustomCode: customCode,
		Status:     httpCode,
	}
}

func (c *Codex) SetErr(err string) *Codex {
	c.Original = fmt.Errorf("%s", err)

	return c
}

func (c *Codex) SetDetail(detail string) *Codex {
	c.Detail = detail

	return c
}

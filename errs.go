package errs

import (
	"fmt"
	"net/http"
	"strings"
)

func write(newErr error, oldErr error, c *Codex) *Format {
	return &Format{
		Prev:     oldErr,
		Original: newErr,
		Msg:      newErr.Error(),
		Trace:    getWithDept(3),
		Codex:    c,
	}
}

func New(str string) *Format {
	return write(fmt.Errorf("%s", str), fmt.Errorf(""), &Codex{})
}

func Chain(err error, a ...any) *Format {
	var s strings.Builder
	for i, v := range a {
		s.WriteString(fmt.Sprintf("%+v", v))
		if i < len(a)-1 {
			s.WriteString(" ")
		}
	}

	return write(fmt.Errorf("%s", s.String()), err, &Codex{})
}

func ChainCodex(err error, c *Codex) *Format {
	return write(c.Original, err, c)
}

func PrintStack(err error) string {
	yes, p := parse(err)
	if !yes || p == nil {
		return "unknown print stack"
	}

	if p.Prev != nil {
		return p.Msg + "\n  " + p.Trace + "\n" + PrintStack(p.Prev)
	}

	return p.Msg + "\n  " + p.Trace
}

type stack struct {
	Msg string `json:"msg"`
	Loc string `json:"loc"`
}

func PrintStackJson(err error) []stack {
	out := []stack{}

	_, p := parse(err)
	if p == nil {
		return out
	}

	out = append(out, stack{
		Msg: p.Msg,
		Loc: p.Trace,
	})

	if p.Prev != nil {
		out = append(out, PrintStackJson(p.Prev)...)
	}

	return out
}

func parse(err error) (bool, *Format) {
	parsed, ok := err.(*Format)
	fmt.Println("parse?", parsed, ok)
	if !ok {
		return false, nil
	}

	return true, parsed
}

func ParseCodex(err error) *Codex {
	defaultCodex := &Codex{
		Title:      "Error unknown",
		Detail:     err.Error(),
		CustomCode: "Codex unknown",
		Status:     http.StatusInternalServerError,
		Original:   err,
	}

	valid, parsed := parse(err)
	if valid && parsed.Codex.Status > 0 {
		return parsed.Codex
	}

	return defaultCodex
}

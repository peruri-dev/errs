package errs

import (
	"fmt"
	"net/http"
	"strings"
)

func write(n error, o error, c *Codex) *Format {
	if c != nil {
		return &Format{
			Prev:     o,
			Original: n,
			Msg:      n.Error(),
			Trace:    getWithDept(3),
			Codex:    c,
		}
	}

	return &Format{
		Prev:     o,
		Original: n,
		Msg:      n.Error(),
		Trace:    getWithDept(3),
	}
}

func New(str string) *Format {
	return write(fmt.Errorf("%s", str), nil, nil)
}

func Chain(err error, a ...any) *Format {
	var s strings.Builder
	for i, v := range a {
		s.WriteString(fmt.Sprintf("%+v", v))
		if i < len(a)-1 {
			s.WriteString(" ")
		}
	}

	return write(fmt.Errorf("%s", s.String()), err, nil)
}

func ChainCodex(err error, c *Codex) *Format {
	return write(c.Original, err, c)
}

func PrintStack(err error) string {
	p := parse(err)
	if p == nil {
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

	p := parse(err)
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

func parse(err error) *Format {
	parsed, ok := err.(*Format)
	if !ok {
		return nil
	}

	return parsed
}

func ParseCodex(err error) *Codex {
	e := parse(err)
	if e == nil {
		return &Codex{
			Title:      "Error unknown",
			Detail:     "Codex unknown",
			CustomCode: "U",
			Status:     http.StatusInternalServerError,
			Original:   err,
		}
	}

	if e.Codex == nil && e.Prev != nil {
		return ParseCodex(e.Prev)
	}

	return e.Codex
}

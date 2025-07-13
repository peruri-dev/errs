package errs

type Format struct {
	Msg      string
	Prev     error
	Original error
	Trace    string
	Codex    *Codex
}

func (f *Format) Error() string {
	if f.Codex.Detail != "" {
		return f.Codex.Detail
	}

	if f.Codex.Title != "" {
		return f.Codex.Title
	}

	return f.Msg
}

type Codex struct {
	Title      string
	Detail     string
	CustomCode string
	Status     int
	Original   error
}

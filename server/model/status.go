package model

type Status int

const (
	Err Status = iota + 1
	Ok
)

func NewStatus(s string) Status {
	if s == "ok" {
		return Ok
	} else if s == "err" {
		return Err
	} else {
		return Err
	}
}

func (s Status) Names() []string {
	return []string{
		"ok",
		"err",
	}
}

func (s Status) String() string {
	return s.Names()[s-1]
}

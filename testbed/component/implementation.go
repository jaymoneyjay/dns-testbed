package component

type Implementation int

const (
	Bind9 Implementation = iota
)

func (i Implementation) String() string {
	return [...]string{"bind9"}[i]
}

package experiment

type Attack int

const (
	Subquery Attack = iota
	SubqueryCNAME
	SubqueryCNAMEQMIN
	SlowDoS
)

func (i Attack) String() string {
	return [...]string{"subquery", "subquery+CNAME", "subquery+CNAME+QMIN", "slowDoS"}[i]
}

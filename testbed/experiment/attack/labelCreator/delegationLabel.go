package labelCreator

type delegationLabel struct {
	*label
}

func NewDelegationLabel() *delegationLabel {
	return &delegationLabel{newLabel("del")}
}

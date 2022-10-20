package labelCreator

type chainLabel struct {
	*label
}

func NewChainLabel() *chainLabel {
	return &chainLabel{newLabel("chain")}
}

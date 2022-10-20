package component

type Component interface {
	Start(implementation Implementation) error
}

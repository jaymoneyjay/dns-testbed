package component

type Logging interface {
	CleanLog() error
	CountQueries() (int, error)
}

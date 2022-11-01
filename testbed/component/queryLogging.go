package component

type QueryLogging interface {
	CleanLog() error
	CountQueries() (int, error)
}

package component

import "path/filepath"

type Resolver struct {
	*Container
	log *log
}

func NewResolver(containerID, buildPath string) *Resolver {
	return &Resolver{
		Container: newContainer(containerID),
		log:       newLog(filepath.Join(buildPath, "logs/query.log")),
	}
}

func (r *Resolver) CleanLog() error {
	return r.log.Clean()
}

func (r *Resolver) CountQueries() (int, error) {
	return r.log.CountQueries()
}

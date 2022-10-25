package component

import (
	"path/filepath"
)

type Resolver struct {
	*Container
	log *containerLog
}

func NewResolver(containerID, buildPath string) (*Resolver, error) {
	container, err := newContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &Resolver{
		Container: container,
		log:       newLog(filepath.Join(buildPath, "logs/query.log"), filepath.Join(buildPath, "logs/general.log")),
	}, nil
}

func (r *Resolver) CleanLog() error {
	return r.log.Clean()
}

func (r *Resolver) CountQueries() (int, error) {
	return r.log.CountQueries()
}

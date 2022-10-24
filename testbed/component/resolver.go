package component

import (
	"errors"
	"fmt"
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
		log:       newLog(filepath.Join(buildPath, "logs/query.containerLog"), filepath.Join(buildPath, "logs/general.containerLog")),
	}, nil
}

func (r *Resolver) Start(implementation Implementation) error {
	switch implementation {
	case Bind9:
		err := r.startBind9()
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("start not implemented for %s", implementation))
	}
	return nil
}

func (r *Resolver) Stop(implementation Implementation) error {
	switch implementation {
	case Bind9:
		err := r.stopBind9()
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("stop not implemented for %s", implementation))
	}
	return nil
}

func (r *Resolver) CleanLog() error {
	return r.log.Clean()
}

func (r *Resolver) CountQueries() (int, error) {
	return r.log.CountQueries()
}

package cqrsexample

import (
	"fmt"

	"github.com/sokool/cqrsexample/events"
	"github.com/sokool/cqrsexample/query"
	"github.com/sokool/gokit/cqrs"
)

type Service struct {
	Query      *query.Query
	Restaurant *Restaurant
}

func NewService() *Service {
	read := query.New()
	write := &Restaurant{
		cqrs.NewRepository(
			factory,
			events.All,
			cqrs.EventHandler(read.Listen)),
	}

	return &Service{
		Query:      read,
		Restaurant: write,
	}
}

type Restaurant struct {
	repository *cqrs.Repository
}

func (s *Restaurant) New() *aggregate {
	return s.repository.Aggregate().(*aggregate)
}

func (s *Restaurant) Load(id string) (*aggregate, error) {
	a, err := s.repository.Load(id)
	if err != nil {
		return nil, err
	}

	r, ok := a.(*aggregate)
	if !ok {
		return nil, fmt.Errorf("wrong aggregate type")
	}

	return r, nil
}

func (s *Restaurant) Save(a *aggregate) error {
	return s.repository.Save(a)
}

func factory() (cqrs.Aggregate, cqrs.DataHandler) {
	r := &aggregate{
		choices: make(map[string]choice),
		menu:    make([]string, 0),
	}
	return r, handler(r)
}

func (a *aggregate) Root() *cqrs.Root {
	return a.root
}

func (a *aggregate) Set(r *cqrs.Root) {
	a.root = r
}

func (a *aggregate) TakeSnapshot() interface{} {
	return nil
}

func (a *aggregate) RestoreSnapshot(s interface{}) error {
	return nil
}

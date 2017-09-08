package query

import (
	"github.com/sokool/cqrsexample/events"
	"github.com/sokool/gokit/cqrs"
)

type Tavern struct {
	ID   int
	UUID string
	Name string
	Info string
	Menu []string
}

type Person struct {
	ID   int
	Name string
}

type Subscriptions struct {
	PersonID int
	TavernID int
}

type Query struct {
	tid     int
	pid     int
	taverns map[string]Tavern
	people  map[string]Person
}

func (q *Query) Listen(a cqrs.CQRSAggregate, ce []cqrs.Event, es []interface{}) {
	for _, event := range es {
		switch e := event.(type) {
		case *events.Created:
			if _, ok := q.taverns[a.String()]; ok {
				break
			}

			q.taverns[a.String()] = Tavern{
				ID:   q.tid,
				UUID: a.ID,
				Name: e.Restaurant,
				Info: e.Info,
				Menu: e.Menu,
			}
			q.tid++
		case *events.MealSelected:
			if _, ok := q.people[e.Person]; ok {
				break
			}

			q.people[e.Person] = Person{
				ID:   q.pid,
				Name: e.Person,
			}
			q.pid++
		}
	}
}

func (q *Query) Taverns() map[string]Tavern {
	return q.taverns
}

func (q *Query) People() map[string]Person {
	return q.people
}

func New() *Query {
	return &Query{
		taverns: map[string]Tavern{},
		people:  map[string]Person{},
	}
}

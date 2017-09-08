package cqrsexample

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sokool/cqrsexample/events"
	"github.com/sokool/gokit/cqrs"
)

type aggregate struct {
	root *cqrs.Root

	name string
	info string
	menu []string

	choices map[string]choice

	created   time.Time
	scheduled time.Time
	//canceled  time.Time
}

type choice struct {
	person string
	meal   string
	on     time.Time
}

func (a *aggregate) Create(name, info string, menu ...string) error {
	if !a.created.IsZero() {
		return fmt.Errorf("restaurant %s is already created", a.name)
	}

	a.root.Apply(&events.Created{
		Restaurant: name,
		Info:       info,
		Menu:       menu,
		At:         time.Now(),
	})

	return nil
}

func (a *aggregate) Schedule(date time.Time) error {
	if a.created.IsZero() {
		return fmt.Errorf("restaurant not created yet")
	}

	if !date.After(time.Now()) {
		return fmt.Errorf("restaurant %s can not be scheduled in past", a.name)
	}

	if len(a.choices) != 0 {
		return fmt.Errorf("can not be rescheduled, food has been chosen by some people")
	}

	if !a.scheduled.IsZero() {
		a.root.Apply(&events.Rescheduled{On: date})
		return nil
	}

	a.root.Apply(&events.Scheduled{On: date})

	return nil
}

func (a *aggregate) ChooseMeal(person, meal string) error {
	if a.created.IsZero() {
		return fmt.Errorf("restaurant is not created yet")
	}

	if a.scheduled.IsZero() {
		return fmt.Errorf("restaurant is not scheduled yet")
	}

	if s, ok := a.choices[person]; ok {
		a.root.Apply(&events.MealChanged{
			Person:       person,
			PreviousMeal: s.meal,
			NewMeal:      meal,
			At:           time.Now()})

		return nil
	}

	a.root.Apply(&events.MealSelected{
		Person: person,
		Meal:   meal,
		At:     time.Now()})

	return nil
}

func handler(a *aggregate) cqrs.DataHandler {
	return func(e interface{}) error {
		switch e := e.(type) {
		case *events.Created:
			a.name, a.info, a.menu = e.Restaurant, e.Info, e.Menu
			a.choices = map[string]choice{}
			a.created = e.At

		case *events.MealSelected:
			a.choices[e.Person] = choice{
				person: e.Person,
				meal:   e.Meal,
				on:     e.At,
			}

		case *events.MealChanged:
			a.choices[e.Person] = choice{
				person: e.Person,
				meal:   e.NewMeal,
				on:     e.At,
			}

		case *events.Scheduled:
			a.scheduled = e.On

		case *events.Rescheduled:
			a.scheduled = e.On

		default:
			return fmt.Errorf("event %s not handled", reflect.TypeOf(e))
		}

		return nil
	}
}

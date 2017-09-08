package events

import "time"

type (
	Created struct {
		Restaurant string
		Info       string
		Menu       []string
		At         time.Time
	}

	MealSelected struct {
		Person string
		Meal   string
		At     time.Time
	}

	MealChanged struct {
		Person       string
		PreviousMeal string
		NewMeal      string
		At           time.Time
	}

	Scheduled struct {
		On time.Time
	}

	Rescheduled struct {
		On time.Time
	}
)

var All = []interface{}{
	&Created{},
	&Scheduled{},
	&Rescheduled{},
	&MealChanged{},
	&MealSelected{},
}

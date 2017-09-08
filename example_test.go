package cqrsexample_test

import (
	"testing"

	"time"

	"github.com/sokool/cqrsexample"
	"github.com/sokool/gokit/test/is"
	"github.com/tonnerre/golang-pretty"
)

var service *cqrsexample.Service = cqrsexample.NewService()

func TestCreateRestaurant(t *testing.T) {
	//WHEN I take new Restaurant aggregate
	restaurant := service.Restaurant.New()

	//THEN I create PasiBus restaurant
	err := restaurant.Create(
		"PasiBus",
		"dobre burgery",
		"BBQ", "EGGY", "Gonzo")

	//I EXPECT no error
	is.NotErr(t, err)

	//THEN I create Zdrowe gary restaurant
	err = restaurant.Create(
		"Zdrowe Gary",
		"description")

	//I EXPECT error already created
	is.Err(t, err, "already created")

}

func TestScheduling(t *testing.T) {
	//WHEN I take new Restaurant aggregate
	restaurant := service.Restaurant.New()

	//THEN I schedule it at +2 days from now.
	err := restaurant.Schedule(time.Now().Add(48 * time.Hour))

	//I EXPECT error restaurant is not created
	is.Err(t, err, "restaurant not created")

	//THEN I Create PasiBus restaurant
	is.NotErr(t, restaurant.Create(
		"PasiBus",
		"dobre burgery",
		"BBQ", "Eggy", "Gonzo"))

	//THEN I schedule it for yesterday.
	err = restaurant.Schedule(time.Now().Add(-24 * time.Hour))

	//I EXPECT error restaurant can not be created in past
	is.Err(t, err, "restaurant can not be created in past")

	//THEN I Schedule PasiBus restaurant at +2 days from now again
	err = restaurant.Schedule(time.Now().Add(2 * 24 * time.Hour))

	//I EXPECT no errors
	is.NotErr(t, err)

	//THEN I Reschedule it at +1 days from now
	err = restaurant.Schedule(time.Now().Add(24 * time.Hour))

	//I EXPECT no errors
	is.NotErr(t, err)

	//THEN I choose Gonzo burger for Tom
	is.NotErr(t, restaurant.ChooseMeal("Tom", "Gonzo"))

	//THEN I Reschedule it at +3 days from now
	err = restaurant.Schedule(time.Now().Add(3 * 24 * time.Hour))

	//I EXPECT food has been chosen by some people error
	is.Err(t, err, "food has been chosen by some people")

}

func TestChooseMeal(t *testing.T) {
	//WHEN I take new Restaurant aggregate

	restaurant := service.Restaurant.New()

	//THEN I choose 'Crazy BBQ' burger for 'Tom'
	err := restaurant.ChooseMeal("Tom", "BBQ")

	//I EXPECT restaurant not created yet error
	is.Err(t, err, "restaurant not created yet")

	//THEN I Create "PasiBurger" restaurant
	is.NotErr(t, restaurant.Create(
		"PasiBus",
		"dobre burgery",
		"BBQ", "Eggy", "Gonzo"))

	//AND I choose 'Crazy BBQ' burger for 'Tom'
	err = restaurant.ChooseMeal("Tom", "BBQ")

	//I EXPECT restaurant is not scheduled yet
	is.Err(t, err, "restaurant not scheduled yet")

	//THEN I Schedule PasiBus restaurant at +5 days from now again
	is.NotErr(t, restaurant.Schedule(time.Now().Add(5*24*time.Hour)))

	//AND I choose 'Crazy BBQ' burger for 'Tom' again
	err = restaurant.ChooseMeal("Tom", "BBQ")

	//I EXPECT no error
	is.NotErr(t, err)

}

func TestScenario(t *testing.T) {

	pasiBus := service.Restaurant.New()
	is.NotErr(t, pasiBus.Create(
		"PasiBus",
		"dobre burgery",
		"BBQ", "Eggy", "Gonzo"))
	is.NotErr(t, pasiBus.Schedule(time.Now().Add(24*time.Hour)))
	is.NotErr(t, pasiBus.ChooseMeal("Tom", "Eggy"))
	is.NotErr(t, pasiBus.ChooseMeal("Greg", "Eggy"))
	is.NotErr(t, pasiBus.ChooseMeal("Tom", "Gonzo"))

	zdroweGary := service.Restaurant.New()
	is.NotErr(t, zdroweGary.Create(
		"Zdrowe Gary",
		"polskie jedzenie",
		"Ogórkowa", "Schabowy", "Pierogi"))
	is.NotErr(t, zdroweGary.Schedule(time.Now().Add(2*24*time.Hour)))
	is.NotErr(t, zdroweGary.Schedule(time.Now().Add(4*24*time.Hour)))
	is.NotErr(t, zdroweGary.ChooseMeal("Cindy", "Schabowy"))
	is.NotErr(t, zdroweGary.ChooseMeal("Tom", "Pierogi"))

	zupapl := service.Restaurant.New()
	is.NotErr(t, zupapl.Create(
		"Zupa.pl",
		"miliardy zup",
		"Ogórkowa", "Pomidorow", "Kalafirowa"))
	is.NotErr(t, zupapl.Schedule(time.Now().Add(3*24*time.Hour)))
	is.NotErr(t, zupapl.ChooseMeal("Joanna", "Pomidorowa"))
	is.NotErr(t, zupapl.ChooseMeal("Tom", "Kalafiorowa"))
	is.NotErr(t, zupapl.ChooseMeal("Cindy", "Pomidorowa"))

	is.NotErr(t, service.Restaurant.Save(pasiBus))
	is.NotErr(t, service.Restaurant.Save(zdroweGary))
	is.NotErr(t, service.Restaurant.Save(zupapl))

	pretty.Println(service.Query.Taverns())
	pretty.Println(service.Query.People())
}

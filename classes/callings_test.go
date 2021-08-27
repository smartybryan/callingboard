package classes

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCallingsFixture(t *testing.T) {
	gunit.Run(new(CallingsFixture), t)
}

type CallingsFixture struct {
	*gunit.Fixture

	callings Callings
}

func (this *CallingsFixture) Setup() {
}

func (this *CallingsFixture) TestDaysInCalling() {
	calling1 := Calling{
		Name:          "Head Honcho",
		Holder:        "User, Joe",
		CustomCalling: false,
		Sustained:     time.Date(2019, 10, 5, 0, 0, 0, 0, time.UTC),
	}
	expected := int(time.Now().Sub(calling1.Sustained).Hours() / 24)

	this.So(calling1.DaysInCalling(), should.Equal, expected)
}

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
	cy, cm, cd := time.Now().Date()
	calling := Calling{
		Name:          "Head Honcho",
		Holder:        "User, Joe",
		CustomCalling: false,
	}

	calling.Sustained = time.Date(cy - 1, cm, cd, 0, 0, 0, 0, time.UTC)
	this.So(calling.DaysInCalling(), should.BeBetween, 363, 368)

	calling.Sustained = time.Date(cy - 1, cm, cd - 5, 0, 0, 0, 0, time.UTC)
	this.So(calling.DaysInCalling(), should.BeBetween, 368, 373)
}

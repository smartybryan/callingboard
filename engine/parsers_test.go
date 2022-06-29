package engine

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestParsersFixture(t *testing.T) {
	gunit.Run(new(ParsersFixture), t)
}

type ParsersFixture struct {
	*gunit.Fixture
}

func (this *ParsersFixture) TestAgeFunctions() {
	birthday := setDate(-20, 0, 0)
	this.So(ageByEndOfYear(birthday), should.Equal, 20)

	birthday = setDate(-20, 0, 2)
	this.So(ageByEndOfYear(birthday), should.Equal, 20)

	birthday = setDate(-18, 0, 2)
	this.So(ageByEndOfYear(birthday), should.Equal, 18)
}

func (this *ParsersFixture) TestCalculateEligibility() {
	this.So(calculateEligibility("2 Jul 2020", false), should.Equal, CallingNotEligible)
	this.So(calculateEligibility("2 Jul 1950", false), should.Equal, CallingNotEligible)
	this.So(calculateEligibility("2 Jul 1950", true), should.Equal, CallingAdult)
	this.So(calculateEligibility("2 Jul 2010", true), should.Equal, CallingYouth)
}

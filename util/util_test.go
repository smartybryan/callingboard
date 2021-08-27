package util

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestUtilFixture(t *testing.T) {
	gunit.Run(new(UtilFixture), t)
}

type UtilFixture struct {
	*gunit.Fixture
}

func (this *UtilFixture) TestPrintableTimeInCalling() {
	this.So(PrintableTimeInCalling(1), should.Equal, "A few days")
	this.So(PrintableTimeInCalling(35), should.Equal, "1 month")
	this.So(PrintableTimeInCalling(90), should.Equal, "3 months")
	this.So(PrintableTimeInCalling(365), should.Equal, "1 year")
	this.So(PrintableTimeInCalling(370), should.Equal, "1 year")
	this.So(PrintableTimeInCalling(720), should.Equal, "1 year, 11 months")
	this.So(PrintableTimeInCalling(730), should.Equal, "2 years")
	this.So(PrintableTimeInCalling(760), should.Equal, "2 years, 1 month")
	this.So(PrintableTimeInCalling(1075), should.Equal, "2 years, 11 months")
	this.So(PrintableTimeInCalling(1095), should.Equal, "3 years")
}

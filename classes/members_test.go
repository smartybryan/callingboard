package classes

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestMembersFixture(t *testing.T) {
	gunit.Run(new(MembersFixture), t)
}

type MembersFixture struct {
	*gunit.Fixture
}

func (this *MembersFixture) Setup() {
}

func (this *MembersFixture) TestAge() {
	cy, cm, cd := time.Now().Date()
	member := Member{
		Name:       "User, Joe",
		Gender:     "M",
		Unbaptized: false,
	}

	member.Birthday = time.Date(cy - 20, cm, cd, 0, 0, 0, 0, time.UTC)
	this.So(member.Age(), should.Equal, 20)

	member.Birthday = time.Date(cy - 20, cm, cd + 2, 0, 0, 0, 0, time.UTC)
	this.So(member.Age(), should.Equal, 19)

	member.Birthday = time.Date(cy - 20, cm, cd + 2, 0, 0, 0, 0, time.UTC)
	this.So(member.AgeByEndOfYear(), should.Equal, 20)
}

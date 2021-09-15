package engine

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestProjectFixture(t *testing.T) {
    gunit.Run(new(ProjectFixture), t)
}

type ProjectFixture struct {
    *gunit.Fixture
}

func (this *ProjectFixture) Setup() {
}

func (this *ProjectFixture) TestUndoRedoTransactions() {
	
}

package exec

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestExec(t *testing.T) { TestingT(t) }

type ExecSuite struct{}

var _ = Suite(&ExecSuite{})

func (s *ExecSuite) TestExec(c *C) {
	o, e := Exec("ls -d /tmp")
	c.Check(o[0], Equals, "/tmp")
	c.Check(e, Equals, nil)

	o, e = Exec("ls -d /tmp/no-file-here")
	c.Check(len(o), Equals, 1)
	c.Check(o[0], Matches, "ls: cannot access '?/tmp/no-file-here'?: No such file or directory")
	c.Check(e, ErrorMatches, "exit status 2")

	o, e = Exec("echo hello")
	c.Check(o[0], Equals, "hello")
	c.Check(e, Equals, nil)

	o, e = Exec("eho hello")
	c.Check(len(o), Equals, 0)
	c.Assert(e, ErrorMatches, ".*executable file not found.*")
}

// +build test

package ninchatstate

func (s *State) TestBreakConnection() func()  { return s.session.TestBreakConnection() }
func (s *State) TestForgetConnection() func() { return s.session.TestForgetConnection() }
func (s *State) TestLoseSession()             { s.session.TestLoseSession() }

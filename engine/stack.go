package engine

type Stack struct {
	X, Y, Z, T float64
}

func (s *Stack) push() {
	s.T = s.Z
	s.Z = s.Y
	s.Y = s.X
}

func (s *Stack) tuck() {
	s.Y = s.Z
	s.Z = s.T
}

func (s *Stack) rollDown() {
	x := s.X
	s.X = s.Y
	s.Y = s.Z
	s.Z = s.T
	s.T = x
}

func (s *Stack) rollUp() {
	x := s.T
	s.T = s.Z
	s.Z = s.Y
	s.Y = s.X
	s.X = x
}

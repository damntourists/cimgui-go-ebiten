package ebitenbackend

type Set [256]bool

func (s *Set) Contains(b byte) bool {
	return s[b]
}

func (s *Set) Add(b byte) {
	s[b] = true
}

func (s *Set) Remove(b byte) {
	s[b] = false
}

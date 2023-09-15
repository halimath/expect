package set

type Set[T comparable] map[T]struct{}

func New[T comparable](vals ...T) Set[T] {
	s := make(Set[T], len(vals))

	for _, val := range vals {
		s.Add(val)
	}

	return s
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) ToSlice() []T {
	ret := make([]T, 0, len(s))

	for v := range s {
		ret = append(ret, v)
	}

	return ret
}

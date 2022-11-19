package expect

type set[T comparable] map[T]struct{}

func NewSet[T comparable](vals ...T) set[T] {
	s := make(set[T], len(vals))

	for _, val := range vals {
		s.Add(val)
	}

	return s
}

func (s set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s set[T]) Remove(v T) {
	delete(s, v)
}

func (s set[T]) ToSlice() []T {
	ret := make([]T, 0, len(s))

	for v := range s {
		ret = append(ret, v)
	}

	return ret
}

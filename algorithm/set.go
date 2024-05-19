package algorithm

type Set[T comparable] map[T]struct{}

func (s Set[T]) Insert(key T) {
	s[key] = struct{}{}
}

func (s Set[T]) Find(key T) bool {
	if _, has := s[key]; has {
		return true
	}
	return false
}

func (s Set[T]) Delete(key T) {
	delete(s, key)
}

func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for k := range s {
		slice = append(slice, k)
	}
	return slice
}

func (s Set[T]) Intersect(v Set[T]) Set[T] {
	a, b := s, v
	if len(b) < len(a) {
		a, b = b, a
	}
	res := Set[T]{}
	for k := range a {
		if b.Find(k) {
			res.Insert(k)
		}
	}
	return res
}

func (s Set[T]) Union(v Set[T]) Set[T] {
	size := len(s)
	if n := len(v); n > size {
		size = n
	}
	res := make(Set[T], size)
	for k := range s {
		res.Insert(k)
	}
	for k := range v {
		res.Insert(k)
	}
	return res
}

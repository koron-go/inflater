/*
Package inflater provides ...
*/
package inflater

import (
	"iter"
)

type Inflater[V any] interface {
	Inflate(seed V) iter.Seq[V]
}

type InflaterFunc[V any] func(seed V) iter.Seq[V]

func (f InflaterFunc[V]) Inflate(seed V) iter.Seq[V] {
	return f(seed)
}

// None provides an Inflater which not inflate anything.
func None[V any]() Inflater[V] {
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {}
	})
}

// Keep provides an Inflater which inflate just only seed.
func Keep[V any]() Inflater[V] {
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			yield(seed)
		}
	})
}

func Map[V any](apply func(V) V) Inflater[V] {
	if apply == nil {
		return Keep[V]()
	}
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			if !yield(apply(seed)) {
				return
			}
		}
	})
}

// Filter provides an Inflater which inflate a seed if check(seed) returns true.
func Filter[V any](check func(V) bool) Inflater[V] {
	if check == nil {
		return Keep[V]()
	}
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			if check(seed) && !yield(seed) {
				return
			}
		}
	})
}

// Parallel2 creates an Inflater with distribute a seed to two Inflaters.
func Parallel2[V any](first, second Inflater[V]) Inflater[V] {
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			for s := range first.Inflate(seed) {
				if !yield(s) {
					return
				}
			}
			for s := range second.Inflate(seed) {
				if !yield(s) {
					return
				}
			}
		}
	})
}

// Parallel creates an Inflater which distibute a seed to multiple Inflaters.
func Parallel[V any](inflaters ...Inflater[V]) Inflater[V] {
	switch len(inflaters) {
	case 0:
		return None[V]()
	case 1:
		return inflaters[0]
	case 2:
		return Parallel2(inflaters[0], inflaters[1])
	default:
		return Parallel2(inflaters[0], Parallel(inflaters[1:]...))
	}
}

// Serial2 creates an Inflater that inflates the result of the first
// Inflater with the second Inflater.
func Serial2[V any](first, second Inflater[V]) Inflater[V] {
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			for s := range first.Inflate(seed) {
				for t := range second.Inflate(s) {
					if !yield(t) {
						return
					}
				}
			}
		}
	})
}

// Serial creates an Inflater that applies multiple Inflaters in sequence to
// its result repeatedly.
func Serial[V any](inflaters ...Inflater[V]) Inflater[V] {
	switch len(inflaters) {
	case 0:
		return None[V]()
	case 1:
		return inflaters[0]
	case 2:
		return Serial2[V](inflaters[0], inflaters[1])
	default:
		return Serial2[V](inflaters[0], Serial(inflaters[1:]...))
	}
}

// Prefix provides an Inflater which inflate with prefixes.
func Prefix[V ~string](prefixes ...V) Inflater[V] {
	if len(prefixes) == 0 {
		return None[V]()
	}
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			for _, prefix := range prefixes {
				if !yield(prefix + seed) {
					return
				}
			}
		}
	})
}

// Suffix provides an Inflater which inflate with suffixes.
func Suffix[V ~string](suffixes ...V) Inflater[V] {
	if len(suffixes) == 0 {
		return None[V]()
	}
	return InflaterFunc[V](func(seed V) iter.Seq[V] {
		return func(yield func(V) bool) {
			for _, suffix := range suffixes {
				if !yield(seed + suffix) {
					return
				}
			}
		}
	})
}

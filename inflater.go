/*
Package `inflater` provides an `Inflater` interface that inflates a value
into multiple values and some operations to combine them `Inflater`s.
*/
package inflater

import (
	"iter"
)

// Inflater is the interface that inflate a value to a sequence of values
type Inflater[V any] interface {
	// Inflate a seed value to a sequence of values.
	Inflate(seed V) iter.Seq[V]
}

// InflateFunc is a wrapper function for using a function as an Inflater.
type InflaterFunc[V any] func(seed V) iter.Seq[V]

// Inflate a value to a sequence of values.
func (f InflaterFunc[V]) Inflate(seed V) iter.Seq[V] {
	return f(seed)
}

// Slice is a wrapper for a slice, which creates an Inflater that returns the
// elements of the slice. This Inflater always ignores the input.
type Slice[V any] []V

// Inflate inflates all elements of the slice. It ignores the input seed.
func (slice Slice[V]) Inflate(seed V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
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

// Map is an Inflater which maps (modify/convert) a value to another value with a function.
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

// Filter provides an Inflater which pass through a seed if check(seed) returns true.
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

// Parallel2 creates an Inflater that inflates one input with two Inflaters and
// concatenates the results into one iter.Seq.
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

// Parallel creates an Inflater that inflates one input with multiple Inflaters
// and concatenates the results into one iter.Seq.
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

// Serial2 creates an Inflater that inflates the input with the first Inflater
// and then inflates the result with the second Inflater.
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

// Serial creates an Inflater that inflates the input with the first Inflater,
// then inflates the result with the second Inflater, and repeats this for all
// the given Inflaters.
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

// Prefix creates an Inflater that prepends each prefix string to one input.
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

// Suffix creates an Inflater that appends each suffix string to the input.
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

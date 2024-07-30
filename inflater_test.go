package inflater_test

import (
	"iter"
	"strings"
	"testing"

	"github.com/koron-go/inflater"
)

type staticInflater []string

func (iter staticInflater) Inflate(string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, s := range iter {
			if !yield(s) {
				return
			}
		}
	}
}

var static1 = staticInflater{"aaa", "bbb", "ccc"}
var static2 = staticInflater{"111", "222", "333"}

func testInflater[T ~string | ~int](t *testing.T, target inflater.Inflater[T], seed T, wants ...T) {
	i := 0
	for got := range target.Inflate(seed) {
		if len(wants) == 0 {
			t.Fatal("length of wants array is insufficient")
		}
		want := wants[0]
		wants = wants[1:]
		if got != want {
			t.Errorf("#%d unmatch: want=%+v got=%+v", i, want, got)
		}
		i++
	}
	if len(wants) > 0 {
		t.Helper()
		t.Errorf("wants array have extra elements")
	}
}

func TestStatic(t *testing.T) {
	testInflater(t, static1, "", "aaa", "bbb", "ccc")
	testInflater(t, static1, "dummy", "aaa", "bbb", "ccc")
	testInflater(t, static2, "", "111", "222", "333")
	testInflater(t, static2, "dummy", "111", "222", "333")
}

func TestPrefix(t *testing.T) {
	testInflater(t, inflater.Prefix("1", "2", "3"), "",
		"1", "2", "3")
	testInflater(t, inflater.Prefix("XXX", "YYY", "ZZZ"), "seed",
		"XXXseed", "YYYseed", "ZZZseed")

	testInflater(t, inflater.Prefix("ONLY"), "", "ONLY")
	testInflater(t, inflater.Prefix("ONLY"), "seed", "ONLYseed")

	// no prefixes, no outputs
	testInflater(t, inflater.Prefix[string](), "")
	testInflater(t, inflater.Prefix[string](), "seed")
}

func TestSuffix(t *testing.T) {
	testInflater(t, inflater.Suffix("1", "2", "3"), "",
		"1", "2", "3")
	testInflater(t, inflater.Suffix("XXX", "YYY", "ZZZ"), "seed",
		"seedXXX", "seedYYY", "seedZZZ")

	testInflater(t, inflater.Suffix("ONLY"), "", "ONLY")
	testInflater(t, inflater.Suffix("ONLY"), "seed", "seedONLY")

	// no suffixes, no outputs
	testInflater(t, inflater.Suffix[string](), "")
	testInflater(t, inflater.Suffix[string](), "seed")
}

func TestSerial2(t *testing.T) {
	testInflater(t, inflater.Serial2(static1, inflater.Suffix("1", "2", "3")), "",
		"aaa1", "aaa2", "aaa3", "bbb1", "bbb2", "bbb3", "ccc1", "ccc2", "ccc3")

	testInflater(t, inflater.Serial2(static1, static2), "",
		"111", "222", "333", "111", "222", "333", "111", "222", "333")
}

func TestParallel2(t *testing.T) {
	testInflater(t, inflater.Parallel2(static1, static2), "",
		"aaa", "bbb", "ccc", "111", "222", "333")
}

func TestNone(t *testing.T) {
	testInflater(t, inflater.None[string](), "")
	testInflater(t, inflater.None[string](), "seed")

	testInflater(t, inflater.Serial2(static1, inflater.None[string]()), "")
}

func TestKeep(t *testing.T) {
	testInflater(t, inflater.Keep[string](), "", "")
	testInflater(t, inflater.Keep[string](), "seed", "seed")

	testInflater(t, inflater.Serial2(static1, inflater.Keep[string]()), "", "aaa", "bbb", "ccc")
}

func TestMap(t *testing.T) {
	testInflater(t, inflater.Serial2(static1, inflater.Map(strings.ToUpper)), "", "AAA", "BBB", "CCC")

	testInflater(t, inflater.Serial2(static1, inflater.Map[string](nil)), "", "aaa", "bbb", "ccc")
}

func TestFilter(t *testing.T) {
	testInflater(t, inflater.Serial2(static1, inflater.Filter(func(s string) bool {
		return s > "aaa"
	})), "", "bbb", "ccc")
	testInflater(t, inflater.Serial2(static1, inflater.Filter(func(s string) bool {
		return s != "bbb"
	})), "", "aaa", "ccc")
	testInflater(t, inflater.Serial2(static1, inflater.Filter(func(s string) bool {
		return s == "bbb"
	})), "", "bbb")

	testInflater(t, inflater.Serial2(static1, inflater.Filter[string](nil)), "", "aaa", "bbb", "ccc")
}

func TestParallel(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		testInflater(t, inflater.Parallel[string](), "")
		testInflater(t, inflater.Parallel[string](), "seed")
	})
	t.Run("1", func(t *testing.T) {
		testInflater(t, inflater.Parallel(static1), "", "aaa", "bbb", "ccc")
		testInflater(t, inflater.Parallel(static2), "", "111", "222", "333")
	})
	t.Run("2", func(t *testing.T) {
		testInflater(t, inflater.Parallel(static1, static2), "", "aaa", "bbb", "ccc", "111", "222", "333")
		testInflater(t, inflater.Parallel(static2, static1), "", "111", "222", "333", "aaa", "bbb", "ccc")
	})
	t.Run("2+", func(t *testing.T) {
		testInflater(t, inflater.Parallel(static1, static2, inflater.Keep[string]()), "seed",
			"aaa", "bbb", "ccc", "111", "222", "333", "seed")
		testInflater(t, inflater.Parallel(static1, static2, inflater.Keep[string](), static1), "seed",
			"aaa", "bbb", "ccc", "111", "222", "333", "seed", "aaa", "bbb", "ccc")
	})
}

func TestSerial(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		testInflater(t, inflater.Serial[string](), "")
		testInflater(t, inflater.Serial[string](), "seed")
	})
	t.Run("1", func(t *testing.T) {
		testInflater(t, inflater.Serial(static1), "", "aaa", "bbb", "ccc")
		testInflater(t, inflater.Serial(static2), "", "111", "222", "333")
	})
	t.Run("2", func(t *testing.T) {
		testInflater(t, inflater.Serial(static1, inflater.Suffix("1", "2", "3")), "",
			"aaa1", "aaa2", "aaa3", "bbb1", "bbb2", "bbb3", "ccc1", "ccc2", "ccc3")
		testInflater(t, inflater.Serial(static1, static2), "",
			"111", "222", "333", "111", "222", "333", "111", "222", "333")
	})
	t.Run("2+", func(t *testing.T) {
		testInflater(t, inflater.Serial(static1, static2, inflater.Keep[string]()), "seed",
			"111", "222", "333", "111", "222", "333", "111", "222", "333")
	})
}

func TestMigemoGraph(t *testing.T) {
	var (
		keep = inflater.Keep[string]()

		roma2hira   = inflater.Suffix(":roma2hira")
		hira2kata   = inflater.Suffix(":hira2kata")
		wide2narrow = inflater.Suffix(":wide2narrow")
		narrow2wide = inflater.Suffix(":narrow2wide")
		dict        = inflater.Suffix(":dict")
	)
	compose := inflater.Parallel(
		keep,
		dict,
		inflater.Serial(roma2hira, inflater.Parallel(
			keep,
			dict,
			inflater.Serial(hira2kata, inflater.Parallel(
				keep,
				dict,
				wide2narrow,
			)),
		)),
		narrow2wide,
	)
	testInflater(t, compose, "query",
		"query",
		"query:dict",
		"query:roma2hira",
		"query:roma2hira:dict",
		"query:roma2hira:hira2kata",
		"query:roma2hira:hira2kata:dict",
		"query:roma2hira:hira2kata:wide2narrow",
		"query:narrow2wide",
	)
}

package inflater_test

import (
	"fmt"

	"github.com/koron-go/inflater"
)

func ExampleSlice() {
	slice := inflater.Slice[string]{"foo", "bar", "baz"}
	for s := range slice.Inflate("IGNORED") {
		fmt.Println(s)
	}
	// Output:
	// foo
	// bar
	// baz
}

func ExamplePrefix() {
	prefix := inflater.Prefix("1st ", "2nd ", "3rd ")
	for s := range prefix.Inflate("item") {
		fmt.Println(s)
	}
	// Output:
	// 1st item
	// 2nd item
	// 3rd item
}

func ExampleSuffix() {
	suffix := inflater.Suffix("-san", "-sama", "-dono")
	for s := range suffix.Inflate("mina") {
		fmt.Println(s)
	}
	// Output:
	// mina-san
	// mina-sama
	// mina-dono
}

func ExampleParallel() {
	parallel := inflater.Parallel(
		inflater.Prefix("1st ", "2nd ", "3rd "),
		inflater.Suffix("-san", "-sama", "-dono"),
	)
	for s := range parallel.Inflate("foo") {
		fmt.Println(s)
	}
	// Output:
	// 1st foo
	// 2nd foo
	// 3rd foo
	// foo-san
	// foo-sama
	// foo-dono
}

func ExampleSerial() {
	serial := inflater.Serial(
		inflater.Prefix("1st ", "2nd ", "3rd "),
		inflater.Suffix("-san", "-sama", "-dono"),
	)
	for s := range serial.Inflate("foo") {
		fmt.Println(s)
	}
	// Output:
	// 1st foo-san
	// 1st foo-sama
	// 1st foo-dono
	// 2nd foo-san
	// 2nd foo-sama
	// 2nd foo-dono
	// 3rd foo-san
	// 3rd foo-sama
	// 3rd foo-dono
}

func ExampleKeep() {
	parallel := inflater.Parallel(
		inflater.Keep[string](),
		inflater.Prefix("1st ", "2nd ", "3rd "),
	)
	for s := range parallel.Inflate("foo") {
		fmt.Println(s)
	}
	// Output:
	// foo
	// 1st foo
	// 2nd foo
	// 3rd foo
}

func ExampleNone() {
	serial := inflater.Serial(
		inflater.Prefix("1st ", "2nd ", "3rd "),
		inflater.Suffix("-san", "-sama", "-dono"),
		inflater.None[string](),
	)
	for s := range serial.Inflate("foo") {
		fmt.Println(s)
	}
	// Output:
}

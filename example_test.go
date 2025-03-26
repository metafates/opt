package opt

import (
	"fmt"
)

func ExampleOpt_IsExplicit() {
	var implicit Opt[string]

	explicit := None[string]()

	fmt.Println(implicit, implicit.IsExplicit())
	fmt.Println(explicit, explicit.IsExplicit())

	// Output:
	// None false
	// None true
}

func ExampleSome() {
	x := Some(2)

	fmt.Println(x)
	// Output: Some(2)
}

func ExampleNone() {
	x := None[int]()

	fmt.Println(x)
	// Output: None
}

func ExampleOpt_IsSome() {
	x := Some(2)
	fmt.Println(x.IsSome())

	x = None[int]()
	fmt.Println(x.IsSome())

	// Output:
	// true
	// false
}

func ExampleOpt_IsSomeAnd() {
	x := Some(2)
	fmt.Println(x.IsSomeAnd(func(x int) bool { return x > 1 }))

	x = Some(0)
	fmt.Println(x.IsSomeAnd(func(x int) bool { return x > 1 }))

	x = None[int]()
	fmt.Println(x.IsSomeAnd(func(x int) bool { return x > 1 }))

	// Output:
	// true
	// false
	// false
}

func ExampleOpt_IsNone() {
	x := Some(2)
	fmt.Println(x.IsNone())

	x = None[int]()
	fmt.Println(x.IsNone())

	// Output:
	// false
	// true
}

func ExampleOpt_IsNoneOr() {
	x := Some(2)
	fmt.Println(x.IsNoneOr(func(x int) bool { return x > 1 }))

	x = Some(0)
	fmt.Println(x.IsNoneOr(func(x int) bool { return x > 1 }))

	x = None[int]()
	fmt.Println(x.IsNoneOr(func(x int) bool { return x > 1 }))

	// Output:
	// true
	// false
	// true
}

func ExampleOpt_TryGet() {
	x := Some(42)
	y := None[int]()

	fmt.Println(x.TryGet())
	fmt.Println(y.TryGet())

	// Output:
	// 42 true
	// 0 false
}

func ExampleOpt_GetOr() {
	fmt.Println(Some("car").GetOr("bike"))
	fmt.Println(None[string]().GetOr("bike"))

	// Output:
	// car
	// bike
}

func ExampleOpt_GetOrElse() {
	k := 10

	fmt.Println(Some(4).GetOrElse(func() int { return 2 * k }))
	fmt.Println(None[int]().GetOrElse(func() int { return 2 * k }))

	// Output:
	// 4
	// 20
}

func ExampleOpt_GetOrEmpty() {
	x := None[int]()
	y := Some(12)

	fmt.Println(x.GetOrEmpty())
	fmt.Println(y.GetOrEmpty())

	// Output:
	// 0
	// 12
}

func ExampleOpt_Map() {
	maybeSomeAge := Some(30)

	maybeSomeYear := maybeSomeAge.Map(func(age int) int { return 2025 - age })

	fmt.Println(maybeSomeYear)

	x := None[int]()

	fmt.Println(x.Map(func(n int) int { return n * 2 }))

	// Output:
	// Some(1995)
	// None
}

func ExampleMap() {
	maybeSomeString := Some("Hello, World!")
	maybeSomeLen := Map(maybeSomeString, func(s string) int { return len(s) })

	fmt.Println(maybeSomeLen)

	x := None[string]()

	fmt.Println(Map(x, func(s string) int { return len(s) }))

	// Output:
	// Some(13)
	// None
}

func ExampleOpt_And() {
	x := Some(2)
	y := None[int]()

	fmt.Println(x.And(y))

	x = None[int]()
	y = Some(42)

	fmt.Println(x.And(y))

	x = Some(2)
	y = Some(42)

	fmt.Println(x.And(y))

	x = None[int]()
	y = None[int]()

	fmt.Println(x.And(y))

	// Output:
	// None
	// None
	// Some(42)
	// None
}

func ExampleOpt_AndThen() {
	div42By := func(divider int) Opt[int] {
		if divider == 0 {
			return None[int]()
		}

		return Some(42 / divider)
	}

	fmt.Println(Some(2).AndThen(div42By))
	fmt.Println(None[int]().AndThen(div42By))
	fmt.Println(Some(0).AndThen(div42By))

	// Output:
	// Some(21)
	// None
	// None
}

func ExampleAndThen() {
	firstRune := func(s string) Opt[rune] {
		runes := []rune(s)

		if len(runes) == 0 {
			return None[rune]()
		}

		return Some(runes[0])
	}

	fmt.Println(AndThen(Some("яблоко"), firstRune))
	fmt.Println(AndThen(Some(""), firstRune))
	fmt.Println(AndThen(None[string](), firstRune))

	// Output:
	// Some(1103)
	// None
	// None
}

func ExampleOpt_Filter() {
	isEven := func(n int) bool { return n%2 == 0 }

	fmt.Println(None[int]().Filter(isEven))
	fmt.Println(Some(3).Filter(isEven))
	fmt.Println(Some(4).Filter(isEven))

	// Output:
	// None
	// None
	// Some(4)
}

func ExampleOpt_Or() {
	x := Some(2)
	y := None[int]()

	fmt.Println(x.Or(y))

	x = None[int]()
	y = Some(100)

	fmt.Println(x.Or(y))

	x = Some(2)
	y = Some(100)

	fmt.Println(x.Or(y))

	x = None[int]()
	y = None[int]()

	fmt.Println(x.Or(y))

	// Output:
	// Some(2)
	// Some(100)
	// Some(2)
	// None
}

func ExampleOpt_OrElse() {
	nobody := func() Opt[string] { return None[string]() }
	vikings := func() Opt[string] { return Some("vikings") }

	fmt.Println(Some("barbarians").OrElse(vikings))
	fmt.Println(None[string]().OrElse(vikings))
	fmt.Println(None[string]().OrElse(nobody))

	// Output:
	// Some(barbarians)
	// Some(vikings)
	// None
}

func ExampleIndexSlice() {
	s := []int{10, 40, 30}

	fmt.Println(IndexSlice(s)(1))
	fmt.Println(IndexSlice(s)(3))
	fmt.Println(Some(0).AndThen(IndexSlice(s)))

	// Output:
	// Some(40)
	// None
	// Some(10)
}

func ExampleIndexMap() {
	m := map[int]int{
		7: 5,
	}

	fmt.Println(IndexMap(m)(7))
	fmt.Println(IndexMap(m)(1))

	fmt.Println(Some(7).AndThen(IndexMap(m)))

	// Output:
	// Some(5)
	// None
	// Some(5)
}

func ExampleFromZero() {
	fmt.Println(FromZero(0))
	fmt.Println(FromZero("foo"))
	fmt.Println(FromZero(""))

	// Output:
	// None
	// Some(foo)
	// None
}

func ExampleFromTuple() {
	foo := func() (int, bool) { return 42, true }
	bar := func() (int, bool) { return 0, false }

	fmt.Println(FromTuple(foo()))

	fmt.Println(FromTuple(bar()))

	// Output:
	// Some(42)
	// None
}

func ExampleFromPtr() {
	value := 42

	var x, y *int

	x = &value
	y = nil

	fmt.Println(FromPtr(x))
	fmt.Println(FromPtr(y))

	// Output:
	// Some(42)
	// None
}

func ExampleOpt_ToPtr() {
	x := Some(42)
	y := None[int]()

	fmt.Println(x.ToPtr() != nil, *x.ToPtr())
	fmt.Println(y.ToPtr())

	// Output:
	// true 42
	// <nil>
}

func ExampleOpt_ToSlice() {
	x := Some(42)
	y := None[int]()
	z := Some([]int{1, 2, 3})

	fmt.Println(x.ToSlice())
	fmt.Println(y.ToSlice())
	fmt.Println(z.ToSlice())

	// Output:
	// [42]
	// []
	// [[1 2 3]]
}

func ExampleOpt_Inspect() {
	x := Some("banana")
	x.Inspect(func(s string) { fmt.Println(s) })

	y := None[string]()
	y.Inspect(func(s string) { fmt.Println(s) })

	// Output: banana
}

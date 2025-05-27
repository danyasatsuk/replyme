package replyme

import (
	"slices"
	"testing"
)

func newFlagInt() *FlagValue[int] {
	return &FlagValue[int]{
		Name:           "test",
		Alias:          "t",
		value:          10,
		preParsedValue: "10",
		hasValue:       true,
	}
}

func newFlagString() *FlagValue[string] {
	return &FlagValue[string]{
		Name:           "test",
		Alias:          "t",
		value:          "test",
		preParsedValue: "test",
		hasValue:       true,
	}
}

func newFlagBool() *FlagValue[bool] {
	return &FlagValue[bool]{
		Name:           "test",
		Alias:          "t",
		value:          true,
		preParsedValue: "true",
		hasValue:       true,
	}
}

func newFlagIntArray() *FlagValue[[]int] {
	return &FlagValue[[]int]{
		Name:           "test",
		Alias:          "t",
		value:          []int{1, 2, 3},
		preParsedValue: "1,2,3",
		hasValue:       true,
	}
}

func newFlagStringArray() *FlagValue[[]string] {
	return &FlagValue[[]string]{
		Name:           "test",
		Alias:          "t",
		value:          []string{"a", "b", "c"},
		preParsedValue: "a,b,c",
		hasValue:       true,
	}
}

func TestFlagValue_ParsedValue(t *testing.T) {
	fInt := newFlagInt()
	fString := newFlagString()
	fBool := newFlagBool()
	fIntArray := newFlagIntArray()
	fStringArray := newFlagStringArray()

	if d, err := fInt.ParsedValue(); err != nil || d != 10 {
		t.Fatalf("failed to parse int: %v", err)
	}
	if d, err := fString.ParsedValue(); err != nil || d != "test" {
		t.Fatalf("failed to parse string: %v", err)
	}
	if d, err := fBool.ParsedValue(); err != nil || d != true {
		t.Fatalf("failed to parse bool: %v", err)
	}
	if d, err := fIntArray.ParsedValue(); err != nil || slices.Equal(d.([]int), []int{1, 2, 3}) != true {
		t.Fatalf("failed to parse int array: %v", err)
	}
	if d, err := fStringArray.ParsedValue(); err != nil || slices.Equal(d.([]string), []string{"a", "b", "c"}) != true {
		t.Fatalf("failed to parse string array: %v", err)
	}
}

func TestFlagValue_ValueType(t *testing.T) {
	fInt := newFlagInt()
	fString := newFlagString()
	fBool := newFlagBool()
	fIntArray := newFlagIntArray()
	fStringArray := newFlagStringArray()

	if d := fInt.ValueType(); d != "int" {
		t.Fatalf("failed to get int type: %v", d)
	}
	if d := fString.ValueType(); d != "string" {
		t.Fatalf("failed to get string type: %v", d)
	}
	if d := fBool.ValueType(); d != "bool" {
		t.Fatalf("failed to get bool type: %v", d)
	}
	if d := fIntArray.ValueType(); d != "[]int" {
		t.Fatalf("failed to get int array type: %v", d)
	}
	if d := fStringArray.ValueType(); d != "[]string" {
		t.Fatalf("failed to get string array type: %v", d)
	}
}

func TestFlagValue_ParseString(t *testing.T) {
	f := &FlagValue[string]{
		Name: "test",
	}
	parse, err := f.Parse("test")
	if err != nil {
		t.Fatal(err)
	}

	if parse != "test" {
		t.Fatalf("failed to parse string: %v", parse)
	}
}

func TestFlagValue_ParseBool(t *testing.T) {
	f := &FlagValue[bool]{
		Name: "test",
	}

	parse, err := f.Parse("true")
	if err != nil {
		t.Fatal(err)
	}

	if parse != true {
		t.Fatalf("failed to parse bool: %v", parse)
	}
}

func TestFlagValue_ParseIntArray(t *testing.T) {
	f := &FlagValue[[]int]{
		Name: "test",
	}

	parse, err := f.Parse("1,2,3")
	if err != nil {
		t.Fatal(err)
	}

	if slices.Equal(parse.([]int), []int{1, 2, 3}) != true {
		t.Fatalf("failed to parse int array: %v", parse)
	}
}

func TestFlagValue_ParseStringArray(t *testing.T) {
	f := &FlagValue[[]string]{
		Name: "test",
	}

	parse, err := f.Parse("a,b,c")
	if err != nil {
		t.Fatal(err)
	}

	if slices.Equal(parse.([]string), []string{"a", "b", "c"}) != true {
		t.Fatalf("failed to parse string array: %v", parse)
	}
}

func newFlags() Flags {
	flags := make([]Flag, 0)
	flags = append(flags, &FlagValue[int]{
		Name:           "inttest",
		preParsedValue: "10",
		value:          10,
		hasValue:       true,
	})
	flags = append(flags, &FlagValue[string]{
		Name:           "stringtest",
		preParsedValue: "test",
		value:          "test",
		hasValue:       true,
	})
	flags = append(flags, &FlagValue[bool]{
		Name:           "booltest",
		preParsedValue: "true",
		value:          true,
		hasValue:       true,
	})
	flags = append(flags, &FlagValue[[]int]{
		Name:           "intarraytest",
		preParsedValue: "1,2,3",
		value:          []int{1, 2, 3},
		hasValue:       true,
	})
	flags = append(flags, &FlagValue[[]string]{
		Name:           "stringarraytest",
		preParsedValue: "a,b,c",
		value:          []string{"a", "b", "c"},
		hasValue:       true,
	})
	return flags
}

func TestFlags_GetFlag(t *testing.T) {
	flags := newFlags()
	i := flags.GetFlagInt("inttest", 0)
	if i != 10 {
		t.Fatalf("failed to get int: %v", i)
	}
	s := flags.GetFlagString("stringtest", "")
	if s != "test" {
		t.Fatalf("failed to get string: %v", s)
	}
	ia := flags.GetFlagIntArray("intarraytest")
	if slices.Equal(ia, []int{1, 2, 3}) != true {
		t.Fatalf("failed to get int array: %v", ia)
	}
	sa := flags.GetFlagStringArray("stringarraytest")
	if slices.Equal(sa, []string{"a", "b", "c"}) != true {
		t.Fatalf("failed to get string array: %v", sa)
	}
}

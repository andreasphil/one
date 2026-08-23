package util_test

import (
	"slices"
	"testing"

	"github.com/andreasphil/one/util"
)

func TestNewSet(t *testing.T) {
	result := util.NewSet[string]()

	if result.Len() != 0 {
		t.Errorf("constructor did not return an empty set")
	}
}

func TestNewSetFrom(t *testing.T) {
	values := []string{"one", "two"}
	result := util.NewSetFrom(values)

	for _, expected := range values {
		if !result.Has(expected) {
			t.Errorf("missing value %s in %v", expected, result)
		}
	}
}

func TestAddOne(t *testing.T) {
	values := []string{"one"}
	set := util.NewSet[string]()
	added := set.Add(values...)

	if added != len(values) {
		t.Errorf("should have added %v value(s), returned %v", len(values), added)
	}

	if !set.Has("one") {
		t.Errorf("missing added value in set: %v", set)
	}
}

func TestAddMultiple(t *testing.T) {
	values := []string{"one", "two", "three"}
	set := util.NewSet[string]()
	added := set.Add(values...)

	if added != len(values) {
		t.Errorf("should have added %v value(s), returned %v", len(values), added)
	}

	if !set.Has("one") || !set.Has("two") || !set.Has("three") {
		t.Errorf("missing added value in set: %v", set)
	}
}

func TestAddDuplicate(t *testing.T) {
	values := []string{"one", "two", "two"}
	set := util.NewSet[string]()
	added := set.Add(values...)

	if added != 2 {
		t.Errorf("should have added %v value(s), returned %v", 2, added)
	}
}

func TestAddExisting(t *testing.T) {
	values := []string{"one", "two", "three"}
	set := util.NewSetFrom([]string{"three"})
	added := set.Add(values...)

	if added != 2 {
		t.Errorf("should have added %v value(s), returned %v", 2, added)
	}
}

func TestDelete(t *testing.T) {
	set := util.NewSetFrom([]string{"one", "two"})
	deleted := set.Delete("one")

	if deleted != 1 {
		t.Errorf("should have deleted value, returned %v", deleted)
	}

	if set.Has("one") {
		t.Errorf("should have deleted value but set still has it: %v", set)
	}
}

func TestNoopDelete(t *testing.T) {
	set := util.NewSetFrom([]string{"one", "two"})
	deleted := set.Delete("three")

	if deleted != 0 {
		t.Errorf("should not have deleted value, returned %v", deleted)
	}

	if set.Len() != 2 {
		t.Errorf("should not have modified set: %v", set)
	}
}

func TestDeleteMultiple(t *testing.T) {
	set := util.NewSetFrom([]string{"one", "two", "three"})
	deleted := set.Delete("one", "three")

	if deleted != 2 {
		t.Errorf("should have deleted values, returned %v", deleted)
	}

	if (set.Has("one")) || set.Has("three") {
		t.Errorf("should have deleted value but set still has it: %v", set)
	}
}

func TestHasReturnsTrue(t *testing.T) {
	set := util.NewSet[string]()
	set.Add("one")
	hasValue := set.Has("one")

	if !hasValue {
		t.Errorf("should have value, returned %v", hasValue)
	}
}

func TestHasReturnsFalse(t *testing.T) {
	set := util.NewSet[string]()
	set.Add("one")
	hasValue := set.Has("two")

	if hasValue {
		t.Errorf("should not have value, returned %v", hasValue)
	}
}

func TestValues(t *testing.T) {
	values := []string{"one", "two"}
	set := util.NewSetFrom(values)
	result := set.Values()

	if len(result) != len(values) {
		t.Errorf("unexpected length: was %v, expected %v", len(result), len(values))
	}

	for _, value := range values {
		if !slices.Contains(result, value) {
			t.Errorf("missing value %v in slice %v", value, result)
		}
	}
}

func TestLen(t *testing.T) {
	set := util.NewSetFrom([]string{"one", "two"})
	length := set.Len()

	if length != 2 {
		t.Errorf("unexpected length %v", length)
	}
}

package util_test

import (
	"slices"
	"testing"

	"github.com/andreasphil/one/util"
)

func TestNewSet(t *testing.T) {
	set := util.NewSet[string]()

	if got := set.Len(); got != 0 {
		t.Errorf("NewSet().Len() = %d, want 0", got)
	}
}

func TestNewSetFrom(t *testing.T) {
	values := []string{"one", "two"}
	set := util.NewSetFrom(values)

	for _, value := range values {
		if !set.Has(value) {
			t.Errorf("NewSetFrom(%v).Has(%q) = false, want true", values, value)
		}
	}
}

func TestSetAdd(t *testing.T) {
	type testcase struct {
		name      string
		initial   []string
		add       []string
		wantAdded int
		wantHas   []string
	}

	testcases := []testcase{
		{
			name:      "one value",
			add:       []string{"one"},
			wantAdded: 1,
			wantHas:   []string{"one"},
		},
		{
			name:      "multiple values",
			add:       []string{"one", "two", "three"},
			wantAdded: 3,
			wantHas:   []string{"one", "two", "three"},
		},
		{
			name:      "duplicate value in the same call",
			add:       []string{"one", "two", "two"},
			wantAdded: 2,
			wantHas:   []string{"one", "two"},
		},
		{
			name:      "value the set already has",
			initial:   []string{"three"},
			add:       []string{"one", "two", "three"},
			wantAdded: 2,
			wantHas:   []string{"one", "two", "three"},
		},
		{
			name:      "no values",
			initial:   []string{"one"},
			add:       nil,
			wantAdded: 0,
			wantHas:   []string{"one"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			set := util.NewSetFrom(tc.initial)

			if got := set.Add(tc.add...); got != tc.wantAdded {
				t.Errorf("Add(%v) = %d, want %d", tc.add, got, tc.wantAdded)
			}

			for _, value := range tc.wantHas {
				if !set.Has(value) {
					t.Errorf("after Add(%v), Has(%q) = false, want true", tc.add, value)
				}
			}

			if got := set.Len(); got != len(tc.wantHas) {
				t.Errorf("after Add(%v), Len() = %d, want %d", tc.add, got, len(tc.wantHas))
			}
		})
	}
}

func TestSetDelete(t *testing.T) {
	type testcase struct {
		name        string
		initial     []string
		delete      []string
		wantDeleted int
		wantRemains []string
	}

	testcases := []testcase{
		{
			name:        "one value",
			initial:     []string{"one", "two"},
			delete:      []string{"one"},
			wantDeleted: 1,
			wantRemains: []string{"two"},
		},
		{
			name:        "multiple values",
			initial:     []string{"one", "two", "three"},
			delete:      []string{"one", "three"},
			wantDeleted: 2,
			wantRemains: []string{"two"},
		},
		{
			name:        "value the set does not have",
			initial:     []string{"one", "two"},
			delete:      []string{"three"},
			wantDeleted: 0,
			wantRemains: []string{"one", "two"},
		},
		{
			name:        "no values",
			initial:     []string{"one", "two"},
			delete:      nil,
			wantDeleted: 0,
			wantRemains: []string{"one", "two"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			set := util.NewSetFrom(tc.initial)

			if got := set.Delete(tc.delete...); got != tc.wantDeleted {
				t.Errorf("Delete(%v) = %d, want %d", tc.delete, got, tc.wantDeleted)
			}

			for _, value := range tc.delete {
				if set.Has(value) {
					t.Errorf("after Delete(%v), Has(%q) = true, want false", tc.delete, value)
				}
			}

			for _, value := range tc.wantRemains {
				if !set.Has(value) {
					t.Errorf("after Delete(%v), Has(%q) = false, want true", tc.delete, value)
				}
			}

			if got := set.Len(); got != len(tc.wantRemains) {
				t.Errorf("after Delete(%v), Len() = %d, want %d", tc.delete, got, len(tc.wantRemains))
			}
		})
	}
}

func TestSetHas(t *testing.T) {
	type testcase struct {
		name    string
		initial []string
		value   string
		want    bool
	}

	testcases := []testcase{
		{
			name:    "value in the set",
			initial: []string{"one"},
			value:   "one",
			want:    true,
		},
		{
			name:    "value not in the set",
			initial: []string{"one"},
			value:   "two",
		},
		{
			name:  "empty set",
			value: "one",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			set := util.NewSetFrom(tc.initial)

			if got := set.Has(tc.value); got != tc.want {
				t.Errorf("NewSetFrom(%v).Has(%q) = %v, want %v", tc.initial, tc.value, got, tc.want)
			}
		})
	}
}

func TestSetEqual(t *testing.T) {
	type testcase struct {
		name string
		a    util.Set[string]
		b    util.Set[string]
		want bool
	}

	testcases := []testcase{
		{
			name: "same values",
			a:    util.NewSetFrom([]string{"one", "two"}),
			b:    util.NewSetFrom([]string{"two", "one"}),
			want: true,
		},
		{
			name: "different values",
			a:    util.NewSetFrom([]string{"one", "two"}),
			b:    util.NewSetFrom([]string{"one", "three"}),
		},
		{
			name: "different lengths",
			a:    util.NewSetFrom([]string{"one", "two"}),
			b:    util.NewSetFrom([]string{"one"}),
		},
		{
			name: "two empty sets",
			a:    util.NewSet[string](),
			b:    util.NewSet[string](),
			want: true,
		},
		{
			name: "empty set and the zero value",
			a:    util.NewSet[string](),
			b:    util.Set[string]{},
			want: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.a.Equal(tc.b); got != tc.want {
				t.Errorf("%v.Equal(%v) = %v, want %v", tc.a.Values(), tc.b.Values(), got, tc.want)
			}

			if got := tc.b.Equal(tc.a); got != tc.want {
				t.Errorf("%v.Equal(%v) = %v, want %v", tc.b.Values(), tc.a.Values(), got, tc.want)
			}
		})
	}
}

func TestSetValues(t *testing.T) {
	values := []string{"one", "two"}
	got := util.NewSetFrom(values).Values()

	if len(got) != len(values) {
		t.Fatalf("Values() = %v, want %d values", got, len(values))
	}

	for _, value := range values {
		if !slices.Contains(got, value) {
			t.Errorf("Values() = %v, want it to contain %q", got, value)
		}
	}
}

func TestSetLen(t *testing.T) {
	if got := util.NewSetFrom([]string{"one", "two"}).Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}

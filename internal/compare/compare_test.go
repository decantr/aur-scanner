package compare

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	left := []string{"foo", "bar", "bar", "baz"}
	right := []string{"qux", "bar", "baz"}
	got := Intersection(left, right)
	want := []string{"bar", "baz"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIntersectionEmpty(t *testing.T) {
	got := Intersection([]string{"foo"}, []string{"bar"})
	if len(got) != 0 {
		t.Fatalf("got %v, want empty slice", got)
	}
}

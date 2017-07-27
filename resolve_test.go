package resolver

import (
	"fmt"
	"testing"
)

func TestResolver(t *testing.T) {
	r, err := New("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	a, err := r.Resolve("google.com")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("a = %+v\n", a)
}

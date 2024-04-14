package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleAllCaps() {
	t := &expect.SpyTB{}
	
	expect.It(t, "HELLO").To(be.AllCaps)
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleAllCaps_fail() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(be.AllCaps)
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to in all caps, but it was not in all caps]
}

func ExampleLen() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(be.Len(be.Eq(5)))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleLen_fail() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(be.Len(be.Eq(4)))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to have length be equal to 4, but it was 5]
}

func ExampleSubstring() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(be.Substring("ell"))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleSubstring_fail() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(be.Substring("goodbye"))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to contain "goodbye"]
}

func Example() {
	t := &expect.SpyTB{}
	
	expect.It(t, "hello").To(
		be.Len(be.Eq(5)),
		be.Eq("hello"),
		be.Substring("ell"),
		be.Not(be.AllCaps),
	)
	
	fmt.Println(t.Result())
	// Output: Test passed
}
func TestStringMatchers(t *testing.T) {
	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It(t, "hello").To(be.Len(be.Eq(5)))
		})
		
		t.Run("failing", func(t *testing.T) {
			spyTB := &expect.SpyTB{}
			expect.It(spyTB, "goodbye").To(be.Len(be.Eq(5)))
			expect.It(t, spyTB).To(spytb.Error("expected goodbye to have length be equal to 5, but it was 7"))
		})
	})
}

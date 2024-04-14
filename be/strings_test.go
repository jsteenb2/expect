package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleHaveAllCaps() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "HELLO").To(be.HaveAllCaps)
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveAllCaps_fail() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(be.HaveAllCaps)
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to in all caps, but it was not in all caps]
}

func ExampleHaveLength() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(be.HaveLength(be.Eq(5)))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveLength_fail() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(be.HaveLength(be.Eq(4)))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to have length be equal to 4, but it was 5]
}

func ExampleHaveSubstring() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(be.HaveSubstring("ell"))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveSubstring_fail() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(be.HaveSubstring("goodbye"))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected hello to contain "goodbye"]
}

func Example() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, "hello").To(
		be.HaveLength(be.Eq(5)),
		be.Eq("hello"),
		be.HaveSubstring("ell"),
		expect.Doesnt(be.HaveAllCaps),
	)
	
	fmt.Println(t.Result())
	// Output: Test passed
}
func TestStringMatchers(t *testing.T) {
	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.Expect(t, "hello").To(be.HaveLength(be.Eq(5)))
		})
		
		t.Run("failing", func(t *testing.T) {
			spyTB := &expect.SpyTB{}
			expect.Expect(spyTB, "goodbye").To(be.HaveLength(be.Eq(5)))
			expect.Expect(t, spyTB).To(spytb.HaveError("expected goodbye to have length be equal to 5, but it was 7"))
		})
	})
}

package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleShallowEq() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.ShallowEq([]string{"hello", "world"}))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleShallowEq_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.ShallowEq([]string{"goodbye", "world"}))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to be equal to [goodbye world]]
}

func ExampleContainingItem() {
	t := &expect.SpyTB{}
	
	anArray := []string{"HELLO", "WORLD"}
	expect.It(t, anArray).To(be.ContainingItem(be.AllCaps))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleContainingItem_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.ContainingItem(be.AllCaps))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to contain an item in all caps, but it did not]
}

func ExampleHaveSize() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.Size[string](be.Eq(2)))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveSize_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.Size[string](be.Eq(3)))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to have a size be equal to 3, but it was 2]
}

func ExampleEveryItem() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.EveryItem(be.Substring("o")))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleEveryItem_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.It(t, anArray).To(be.EveryItem(be.Substring("h")))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to have every item contain "h"]
}

func TestArrayMatchers(t *testing.T) {
	t.Run("contain with other matcher to find matcher", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It(t, []string{"hello", "WORLD"}).To(be.ContainingItem(be.AllCaps))
		})
		
		t.Run("failing", func(t *testing.T) {
			t.Run("equal to", func(t *testing.T) {
				spytb.VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					be.ContainingItem(be.Eq("goodbye")),
					`expected [hello world] to contain an item be equal to "goodbye", but it did not`,
				)
			})
			t.Run("all caps", func(t *testing.T) {
				spytb.VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					be.ContainingItem(be.AllCaps),
					`expected [hello world] to contain an item in all caps, but it did not`,
				)
			})
		})
	})
}

package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleShallowEquals() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.ShallowEquals([]string{"hello", "world"}))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleShallowEquals_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.ShallowEquals([]string{"goodbye", "world"}))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to be equal to [goodbye world]]
}

func ExampleContainItem() {
	t := &expect.SpyTB{}
	
	anArray := []string{"HELLO", "WORLD"}
	expect.Expect(t, anArray).To(be.ContainItem(be.HaveAllCaps))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleContainItem_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.ContainItem(be.HaveAllCaps))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to contain an item in all caps, but it did not]
}

func ExampleHaveSize() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.HaveSize[string](be.Eq(2)))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveSize_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.HaveSize[string](be.Eq(3)))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to have a size be equal to 3, but it was 2]
}

func ExampleEveryItem() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.EveryItem(be.HaveSubstring("o")))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleEveryItem_fail() {
	t := &expect.SpyTB{}
	
	anArray := []string{"hello", "world"}
	expect.Expect(t, anArray).To(be.EveryItem(be.HaveSubstring("h")))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected [hello world] to have every item contain "h"]
}

func TestArrayMatchers(t *testing.T) {
	t.Run("contain with other matcher to find matcher", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.Expect(t, []string{"hello", "WORLD"}).To(be.ContainItem(be.HaveAllCaps))
		})
		
		t.Run("failing", func(t *testing.T) {
			t.Run("equal to", func(t *testing.T) {
				spytb.VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					be.ContainItem(be.Eq("goodbye")),
					`expected [hello world] to contain an item be equal to "goodbye", but it did not`,
				)
			})
			t.Run("all caps", func(t *testing.T) {
				spytb.VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					be.ContainItem(be.HaveAllCaps),
					`expected [hello world] to contain an item in all caps, but it did not`,
				)
			})
		})
	})
}

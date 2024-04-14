package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleEq() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Eq(5))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleEq_fail() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Eq(4))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be equal to 4, but it was 5]
}

func ExampleGreater() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Greater(4))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleGreater_fail() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Greater(6))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be greater than 6, but it was 5]
}

func ExampleLess() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Less(6))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleLess_fail() {
	t := &expect.SpyTB{}
	expect.It(t, 5).To(be.Less(4))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be less than 4, but it was 5]
}

func TestComparisonMatchers(t *testing.T) {
	t.Run("Less than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It(t, 5).To(be.Less(6))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, be.Less(6), "expected 6 to be less than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 6, be.Less(3), "expected 6 to be less than 3, but it was 6")
		})
	})
	
	t.Run("Greater than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It(t, 5).To(be.Greater(4))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, be.Greater(6), "expected 6 to be greater than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 2, be.Greater(10), "expected 2 to be greater than 10, but it was 2")
		})
	})
	
	t.Run("equal to with empty strings", func(t *testing.T) {
		t.Run("when it is an empty string, failing output should be quoted", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				"",
				be.Eq("Bob"),
				`expected "" to be equal to "Bob", but it was ""`,
			)
		})
	})
}

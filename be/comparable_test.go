package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleEqual() {
	t := &expect.SpyTB{}
	expect.Expect(t, 2).To(be.Eq(2))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleEqual_fail() {
	t := &expect.SpyTB{}
	expect.Expect(t, 2).To(be.Eq(1))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 2 to be equal to 1, but it was 2]
}

func ExampleEqualTo() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.Eq(5))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleEqualTo_fail() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.Eq(4))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be equal to 4, but it was 5]
}

func ExampleGreaterThan() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.GreaterThan(4))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleGreaterThan_fail() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.GreaterThan(6))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be greater than 6, but it was 5]
}

func ExampleLessThan() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.LessThan(6))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleLessThan_fail() {
	t := &expect.SpyTB{}
	expect.Expect(t, 5).To(be.LessThan(4))
	fmt.Println(t.Result())
	// Output: Test failed: [expected 5 to be less than 4, but it was 5]
}

func TestComparisonMatchers(t *testing.T) {
	t.Run("Less than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.Expect(t, 5).To(be.LessThan(6))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, be.LessThan(6), "expected 6 to be less than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 6, be.LessThan(3), "expected 6 to be less than 3, but it was 6")
		})
	})
	
	t.Run("Greater than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.Expect(t, 5).To(be.GreaterThan(4))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, be.GreaterThan(6), "expected 6 to be greater than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 2, be.GreaterThan(10), "expected 2 to be greater than 10, but it was 2")
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

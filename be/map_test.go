package be_test

import (
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleWithAnyValue() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, map[string]string{"hello": "world"}).To(be.HaveKey("goodbye", be.WithAnyValue[string]()))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected map[hello:world] to have key goodbye, but it did not]
}

func ExampleHaveKey_fail() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, map[string]int{"score": 4}).To(be.HaveKey("score", be.GreaterThan(5).And(be.LessThan(10))))
	
	fmt.Println(t.Result())
	// Output: Test failed: [expected map[score:4] to have key score with value be greater than 5 and be less than 10, but it was 4]
}

func ExampleHaveKey() {
	t := &expect.SpyTB{}
	
	expect.Expect(t, map[string]string{"hello": "world"}).To(be.HaveKey("hello", be.Eq("world")))
	
	fmt.Println(t.Result())
	// Output: Test passed
}

func TestMapMatching(t *testing.T) {
	t.Run("HasKey WithValue", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.Expect(t, map[string]string{"hello": "world"}).To(be.HaveKey("hello", be.Eq("world")))
			expect.Expect(t, map[string]int{"score": 7}).To(be.HaveKey("score", be.GreaterThan(5).And(be.LessThan(10))))
		})
	})
	
	t.Run("failures", func(t *testing.T) {
		t.Run("missing key", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				map[string]string{"hello": "world"},
				be.HaveKey("goodbye", be.WithAnyValue[string]()),
				`expected map[hello:world] to have key goodbye, but it did not`,
			)
		})
		
		t.Run("key exists but value does not match", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				map[string]string{"hello": "world"},
				be.HaveKey("hello", be.Eq("goodbye")),
				`expected map[hello:world] to have key hello with value be equal to "goodbye", but it was "world"`,
			)
		})
	})
}

package expect_test

import (
	"errors"
	"fmt"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
)

func ExampleInspector_To() {
	t := &expect.SpyTB{}
	expect.It(t, "Pepper").To(be.Eq("Stanley"))
	fmt.Println(t.Result())
	// Output: Test failed: [expected "Pepper" to be equal to "Stanley", but it was "Pepper"]
}

func ExampleMatcher_Or() {
	t := &expect.SpyTB{}
	tshirt := TShirt{Colour: "yellow"}
	
	expect.It(t, tshirt).To(HaveColour("blue").Or(HaveColour("red")))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the t-shirt to have colour "blue" or have colour "red", but it was "yellow"]
}

func ExampleNot() {
	t := &expect.SpyTB{}
	
	tshirt := TShirt{Colour: "yellow"}
	
	expect.It(t, tshirt).To(be.Not(HaveColour("yellow")))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the t-shirt to not have colour "yellow"]
}

func ExampleMatcher_And() {
	t := &expect.SpyTB{}
	player := Player{
		Name:   "Chris",
		Points: 11,
	}
	
	expect.It(t, player).To(HaveScore(
		be.Greater(5).And(be.Less(10)),
	))
	fmt.Println(t.Result())
	// Output: Test failed: [expected Player Chris to score be greater than 5 and be less than 10, but it was 11]
}

func ExampleExpectNoError() {
	t := &expect.SpyTB{}
	
	err := errors.New("oh no")
	
	expect.NoError(t, err)
	fmt.Println(t.Result())
	// Output: Test failed: [unexpected error: oh no]
}

func ExampleExpectError() {
	t := &expect.SpyTB{}
	
	err := errors.New("oh no")
	
	expect.Error(t, err)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleExpectErrorOfType() {
	t := &expect.SpyTB{}
	
	unauthorised := errors.New("unauthorised")
	wrappedErr := fmt.Errorf("oh no: %w", unauthorised)
	
	expect.ErrorOfType(t, wrappedErr, unauthorised)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleExpectErrorOfType_failing() {
	t := &expect.SpyTB{}
	
	unauthorised := errors.New("unauthorised")
	wrappedErr := fmt.Errorf("oh no: %w", unauthorised)
	
	expect.ErrorOfType(t, wrappedErr, errors.New("not found"))
	fmt.Println(t.Result())
	// Output: Test failed: [expected error of type *errors.errorString, but got "oh no: unauthorised"]
}

func TestMatching(t *testing.T) {
	t.Run("passing example", func(t *testing.T) {
		expect.It(t, "hello").To(
			be.Len(be.Eq(5)),
			be.Eq("hello"),
			be.Substring("ell"),
			be.Not(be.AllCaps),
		)
	})
	
	t.Run("combining failures", func(t *testing.T) {
		t.Run("when it has a but and both failed", func(t *testing.T) {
			someString := "goodbye"
			result1 := be.Len(be.Eq(5))(someString)
			result2 := be.AllCaps(someString)
			
			expected := expect.MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was 7 and it was not in all caps",
			}
			
			actual := result1.Combine(result2)
			expect.It(t, actual).To(be.Eq(expected))
		})
		
		t.Run("when nothing fails", func(t *testing.T) {
			someString := "HELLO"
			result1 := be.Len(be.Eq(5))(someString)
			result2 := be.AllCaps(someString)
			
			expected := expect.MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     true,
			}
			
			actual := result1.Combine(result2)
			expect.It(t, actual).To(be.Eq(expected))
		})
		
		t.Run("when first match is passing but second is failing", func(t *testing.T) {
			someString := "hello"
			result1 := be.Len(be.Eq(5))(someString)
			result2 := be.AllCaps(someString)
			
			expected := expect.MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was not in all caps",
			}
			
			actual := result1.Combine(result2)
			expect.It(t, actual).To(be.Eq(expected))
		})
		
		t.Run("when first match is failing but second is passing", func(t *testing.T) {
			someString := "GOODBYE"
			result1 := be.Len(be.Eq(5))(someString)
			result2 := be.AllCaps(someString)
			
			expected := expect.MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was 7",
			}
			
			actual := result1.Combine(result2)
			expect.It(t, actual).To(be.Eq(expected))
		})
	})
}

type TShirt struct {
	Colour string
}

func (t TShirt) String() string {
	return "the t-shirt"
}

func HaveColour(colour string) expect.Matcher[TShirt] {
	return func(t TShirt) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have colour %q", colour),
			Matches:     t.Colour == colour,
			But:         fmt.Sprintf("it was %q", t.Colour),
		}
	}
}

type Player struct {
	Name   string
	Points int
}

func (s Player) String() string {
	return fmt.Sprintf("Player %s", s.Name)
}

func HaveScore(matcher expect.Matcher[int]) expect.Matcher[Player] {
	return func(s Player) expect.MatchResult {
		result := matcher(s.Points)
		return expect.MatchResult{
			Description: "score " + result.Description,
			Matches:     result.Matches,
			But:         fmt.Sprintf("it was %d", s.Points),
		}
	}
}

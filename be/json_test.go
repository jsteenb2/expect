package be_test

import (
	"bytes"
	"fmt"
	"io"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
)

func ExampleParse() {
	t := &expect.SpyTB{}
	
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	HasName := func(want string) expect.Matcher[Person] {
		return func(p Person) expect.MatchResult {
			return expect.MatchResult{
				Description: fmt.Sprintf("name is %s", p.Name),
				Matches:     p.Name == want,
				But:         fmt.Sprintf("name is %s", want),
				SubjectName: "Person",
			}
		}
	}
	
	someJSON := bytes.NewBuffer([]byte(`{"name": "John", "age": 42}`))
	
	expect.Expect[io.Reader](t, someJSON).To(be.Parse[Person](HasName("John")))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleParse_fail() {
	t := &expect.SpyTB{}
	
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	someJSON := bytes.NewBuffer([]byte(`invalid json`))
	
	expect.Expect[io.Reader](t, someJSON).To(be.Parse[Person](be.Eq(Person{
		Name: "Pepper",
		Age:  14,
	})))
	fmt.Println(t.Result())
	// Output: Test failed: [expected JSON to be parseable into be_test.Person, but it could not be parsed: invalid character 'i' looking for beginning of value]
}

package bejson_test

import (
	"bytes"
	"fmt"
	"io"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/be/bejson"
)

func ExampleParsed() {
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
	
	expect.It[io.Reader](t, someJSON).To(bejson.Parsed[Person](HasName("John")))
	fmt.Printf("%s\n", t)
	// Output: Test passed
}

func ExampleParsed_fail() {
	t := &expect.SpyTB{}
	
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	someJSON := bytes.NewBuffer([]byte(`invalid json`))
	
	expect.It[io.Reader](t, someJSON).To(bejson.Parsed[Person](be.Eq(Person{
		Name: "Pepper",
		Age:  14,
	})))
	fmt.Printf("%s\n", t)
	// Output: Test failed: [expected JSON to be parseable into bejson_test.Person, but it could not be parsed: invalid character 'i' looking for beginning of value]
}

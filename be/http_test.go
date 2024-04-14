package be_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleBeOK() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusOK)
	
	expect.Expect(t, res.Result()).To(be.BeOK)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleBeOK_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusNotFound)
	
	expect.Expect(t, res.Result()).To(be.BeOK)
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have status of 200, but it was 404]
}

func ExampleHaveBody() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Body.WriteString("Hello, world")
	
	expect.Expect(t, res.Result()).To(be.HaveBody(be.HaveString(be.Eq("Hello, world"))))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveBody_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Body.WriteString("Hello, world")
	
	expect.Expect(t, res.Result()).To(be.HaveBody(be.HaveString(be.Eq("Goodbye, world"))))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response body to be equal to "Goodbye, world", but it was "Hello, world"]
}

func ExampleHaveHeader() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/html")
	
	expect.Expect(t, res.Result()).To(be.HaveHeader("Content-Type", "text/html"))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveJSONHeader() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "application/json")
	
	expect.Expect(t, res.Result()).To(be.HaveJSONHeader)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveHeader_multiple() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Encoding", "gzip")
	res.Header().Add("Content-Type", "text/html")
	
	expect.Expect(t, res.Result()).To(
		be.HaveHeader("Content-Encoding", "gzip"),
		be.HaveHeader("Content-Type", "text/html"),
	)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveHeader_multiple_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/xml")
	
	expect.Expect(t, res.Result()).To(
		be.HaveHeader("Content-Encoding", "gzip"),
		be.HaveHeader("Content-Type", "text/html"),
	)
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have header "Content-Encoding" of "gzip", but it was "" expected the response to have header "Content-Type" of "text/html", but it was "text/xml"]
}

func ExampleHaveStatus() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)
	
	expect.Expect(t, res.Result()).To(be.HaveStatus(http.StatusTeapot))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHaveStatus_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)
	
	expect.Expect(t, res.Result()).To(be.HaveStatus(http.StatusNotFound))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have status of 404, but it was 418]
}

func ExampleHaveHeader_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/xml")
	
	expect.Expect(t, res.Result()).To(be.HaveHeader("Content-Type", "text/html"))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have header "Content-Type" of "text/html", but it was "text/xml"]
}

func TestHTTPTestMatchers(t *testing.T) {
	t.Run("Body matching", func(t *testing.T) {
		t.Run("simple string match", func(t *testing.T) {
			res := httptest.NewRecorder()
			
			res.Body.WriteString("Hello, world")
			
			// see how we can compose matchers together!
			expect.Expect(t, res.Result()).To(be.HaveBody(be.HaveString(be.Eq("Hello, world"))))
		})
		
		t.Run("simple string mismatch", func(t *testing.T) {
			res := httptest.NewRecorder()
			
			res.Body.WriteString("Hello, world")
			
			spytb.VerifyFailingMatcher(
				t,
				res.Result(),
				be.HaveBody(be.HaveString(be.Eq("Goodbye, world"))),
				`expected the response body to be equal to "Goodbye, world", but it was "Hello, world"`,
			)
		})
		
		t.Run("failing to read body", func(t *testing.T) {
			res := httptest.NewRecorder().Result()
			res.Body = FailingIOReadCloser{Error: fmt.Errorf("oops")}
			
			spytb.VerifyFailingMatcher(
				t,
				res,
				be.HaveBody(be.HaveString(be.Eq("Goodbye, world"))),
				"expected the response body to have data in io.Reader, but it could not be read",
			)
		})
		
		t.Run("example of matching JSON", func(t *testing.T) {
			type Todo struct {
				Name        string    `json:"name"`
				Completed   bool      `json:"completed"`
				LastUpdated time.Time `json:"last_updated"`
			}
			
			WithCompletedTODO := func(todo Todo) expect.MatchResult {
				return expect.MatchResult{
					Description: "have a completed todo",
					Matches:     todo.Completed,
					But:         "it wasn't complete",
				}
			}
			WithTodoNameOf := func(todoName string) expect.Matcher[Todo] {
				return func(todo Todo) expect.MatchResult {
					return expect.MatchResult{
						Description: fmt.Sprintf("have a todo name of %q", todoName),
						Matches:     todo.Name == todoName,
						But:         fmt.Sprintf("it was %q", todo.Name),
					}
				}
			}
			
			t.Run("with completed todo", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": true}`)
				expect.Expect(t, res.Result()).To(be.HaveBody(be.Parse[Todo](WithCompletedTODO)))
			})
			
			t.Run("with incomplete todo", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					be.HaveBody(be.Parse[Todo](WithCompletedTODO)),
					"expected the response body to have a completed todo, but it wasn't complete",
				)
			})
			
			t.Run("with a todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				expect.Expect(t, res.Result()).To(be.HaveBody(be.Parse[Todo](WithTodoNameOf("Finish the side project"))))
			})
			
			t.Run("with incorrect todo name and not completed", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					be.HaveBody(be.Parse[Todo](WithTodoNameOf("Bacon").And(WithCompletedTODO))),
					`expected the response body to have a todo name of "Bacon" and have a completed todo, but it was "Egg" and it wasn't complete`,
				)
			})
			
			t.Run("with incorrect todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					be.HaveBody(be.Parse[Todo](WithTodoNameOf("Bacon"))),
					`expected the response body to have a todo name of "Bacon", but it was "Egg"`,
				)
			})
			
			t.Run("compose the matchers", func(t *testing.T) {
				res := httptest.NewRecorder()
				
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				res.Header().Add("content-type", "application/json")
				
				expect.Expect(t, res.Result()).To(
					be.BeOK,
					be.HaveJSONHeader,
					be.HaveBody(be.Parse[Todo](WithTodoNameOf("Egg").And(expect.Not(WithCompletedTODO)))),
				)
			})
			
			t.Run("compose the matchers to show an unexpected response", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &expect.SpyTB{}
				
				res.Body.WriteString(`{"name": "Egg", "completed": true}`)
				
				expect.Expect(spyTB, res.Result()).To(
					be.BeOK,
					be.HaveJSONHeader,
					be.HaveBody(be.Parse[Todo](WithTodoNameOf("Egg").And(expect.Not(WithCompletedTODO)))),
				)
				expect.Expect(t, spyTB).To(spytb.HaveError(`expected the response body to have a todo name of "Egg" and not have a completed todo`))
				expect.Expect(t, spyTB).To(spytb.HaveError(`expected the response to have header "content-type" of "application/json", but it was ""`))
			})
		})
	})
	
	t.Run("Status code matchers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			t.Run("positive happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusOK)
				expect.Expect(t, res.Result()).To(be.BeOK)
			})
			
			t.Run("negation on happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusTeapot)
				expect.Expect(t, res.Result()).To(expect.Not(be.BeOK))
			})
			
			t.Run("failure message", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &expect.SpyTB{}
				res.WriteHeader(http.StatusNotFound)
				expect.Expect(spyTB, res.Result()).To(be.BeOK)
				expect.Expect(t, spyTB).To(spytb.HaveError(`expected the response to have status of 200, but it was 404`))
			})
		})
		
		t.Run("user defined status", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.WriteHeader(http.StatusTeapot)
			expect.Expect(t, res.Result()).To(be.HaveStatus(http.StatusTeapot))
		})
	})
	
	t.Run("Header matchers", func(t *testing.T) {
		t.Run("happy path multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.Header().Add("Content-Encoding", "gzip")
			res.Header().Add("Content-Type", "text/html")
			
			expect.Expect(t, res.Result()).To(
				be.HaveHeader("Content-Encoding", "gzip"),
				be.HaveHeader("Content-Type", "text/html"),
			)
		})
		
		t.Run("unhappy path with multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &expect.SpyTB{}
			res.Header().Add("Content-Type", "text/xml")
			
			expect.Expect(spyTB, res.Result()).To(
				be.HaveHeader("Content-Encoding", "gzip"),
				be.HaveHeader("Content-Type", "text/html"),
			)
			expect.Expect(t, spyTB).To(
				spytb.HaveError(`expected the response to have header "Content-Encoding" of "gzip", but it was ""`),
				spytb.HaveError(`expected the response to have header "Content-Type" of "text/html", but it was "text/xml"`),
			)
		})
	})
}

type FailingIOReadCloser struct {
	Error error
}

func (f FailingIOReadCloser) Read(p []byte) (n int, err error) {
	return 0, f.Error
}

func (f FailingIOReadCloser) Close() error {
	return nil
}

package behttp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/be/behttp"
	"github.com/jsteenb2/expect/be/beio"
	"github.com/jsteenb2/expect/be/bejson"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleRespBody() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Body.WriteString("Hello, world")
	
	expect.It(t, res.Result()).To(behttp.RespBody(beio.String(be.Eq("Hello, world"))))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleRespBody_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Body.WriteString("Hello, world")
	
	expect.It(t, res.Result()).To(behttp.RespBody(beio.String(be.Eq("Goodbye, world"))))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response body to be equal to "Goodbye, world", but it was "Hello, world"]
}

func ExampleHeader() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/html")
	
	expect.It(t, res.Result()).To(behttp.Header("Content-Type", "text/html"))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleContentTypeJSON() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "application/json")
	
	expect.It(t, res.Result()).To(behttp.ContentTypeJSON)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHeader_multiple() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Encoding", "gzip")
	res.Header().Add("Content-Type", "text/html")
	
	expect.It(t, res.Result()).To(
		behttp.Header("Content-Encoding", "gzip"),
		behttp.Header("Content-Type", "text/html"),
	)
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleHeader_multiple_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/xml")
	
	expect.It(t, res.Result()).To(
		behttp.Header("Content-Encoding", "gzip"),
		behttp.Header("Content-Type", "text/html"),
	)
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have header "Content-Encoding" of "gzip", but it was "" expected the response to have header "Content-Type" of "text/html", but it was "text/xml"]
}

func ExampleStatus() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)
	
	expect.It(t, res.Result()).To(behttp.Status(http.StatusTeapot))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleStatus_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)
	
	expect.It(t, res.Result()).To(behttp.Status(http.StatusNotFound))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have status of 404, but it was 418]
}

func ExampleHeader_fail() {
	t := &expect.SpyTB{}
	res := httptest.NewRecorder()
	res.Header().Add("Content-Type", "text/xml")
	
	expect.It(t, res.Result()).To(behttp.Header("Content-Type", "text/html"))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the response to have header "Content-Type" of "text/html", but it was "text/xml"]
}

func TestHTTPTestMatchers(t *testing.T) {
	t.Run("Body matching", func(t *testing.T) {
		t.Run("simple string match", func(t *testing.T) {
			res := httptest.NewRecorder()
			
			res.Body.WriteString("Hello, world")
			
			// see how we can compose matchers together!
			expect.It(t, res.Result()).To(behttp.RespBody(beio.String(be.Eq("Hello, world"))))
		})
		
		t.Run("simple string mismatch", func(t *testing.T) {
			res := httptest.NewRecorder()
			
			res.Body.WriteString("Hello, world")
			
			spytb.VerifyFailingMatcher(
				t,
				res.Result(),
				behttp.RespBody(beio.String(be.Eq("Goodbye, world"))),
				`expected the response body to be equal to "Goodbye, world", but it was "Hello, world"`,
			)
		})
		
		t.Run("failing to read body", func(t *testing.T) {
			res := httptest.NewRecorder().Result()
			res.Body = FailingIOReadCloser{Error: fmt.Errorf("oops")}
			
			spytb.VerifyFailingMatcher(
				t,
				res,
				behttp.RespBody(beio.String(be.Eq("Goodbye, world"))),
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
				expect.It(t, res.Result()).To(behttp.RespBody(bejson.Parsed[Todo](WithCompletedTODO)))
			})
			
			t.Run("with incomplete todo", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					behttp.RespBody(bejson.Parsed[Todo](WithCompletedTODO)),
					"expected the response body to have a completed todo, but it wasn't complete",
				)
			})
			
			t.Run("with a todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				expect.It(t, res.Result()).To(behttp.RespBody(bejson.Parsed[Todo](WithTodoNameOf("Finish the side project"))))
			})
			
			t.Run("with incorrect todo name and not completed", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					behttp.RespBody(bejson.Parsed[Todo](WithTodoNameOf("Bacon").And(WithCompletedTODO))),
					`expected the response body to have a todo name of "Bacon" and have a completed todo, but it was "Egg" and it wasn't complete`,
				)
			})
			
			t.Run("with incorrect todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				
				spytb.VerifyFailingMatcher(
					t,
					res.Result(),
					behttp.RespBody(bejson.Parsed[Todo](WithTodoNameOf("Bacon"))),
					`expected the response body to have a todo name of "Bacon", but it was "Egg"`,
				)
			})
			
			t.Run("compose the matchers", func(t *testing.T) {
				res := httptest.NewRecorder()
				
				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				res.Header().Add("content-type", "application/json")
				
				expect.It(t, res.Result()).To(
					behttp.Status(http.StatusOK),
					behttp.ContentTypeJSON,
					behttp.RespBody(bejson.Parsed[Todo](WithTodoNameOf("Egg").And(be.Not(WithCompletedTODO)))),
				)
			})
			
			t.Run("compose the matchers to show an unexpected response", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &expect.SpyTB{}
				
				res.Body.WriteString(`{"name": "Egg", "completed": true}`)
				
				expect.It(spyTB, res.Result()).To(
					behttp.Status(http.StatusOK),
					behttp.ContentTypeJSON,
					behttp.RespBody(bejson.Parsed[Todo](WithTodoNameOf("Egg").And(be.Not(WithCompletedTODO)))),
				)
				expect.It(t, spyTB).To(spytb.Error(`expected the response body to have a todo name of "Egg" and not have a completed todo`))
				expect.It(t, spyTB).To(spytb.Error(`expected the response to have header "Content-Type" of "application/json", but it was ""`))
			})
		})
	})
	
	t.Run("Status code matchers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			t.Run("positive happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusOK)
				expect.It(t, res.Result()).To(behttp.Status(http.StatusOK))
			})
			
			t.Run("negation on happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusTeapot)
				expect.It(t, res.Result()).To(be.Not(behttp.Status(http.StatusOK)))
			})
			
			t.Run("failure message", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &expect.SpyTB{}
				res.WriteHeader(http.StatusNotFound)
				expect.It(spyTB, res.Result()).To(behttp.Status(http.StatusOK))
				expect.It(t, spyTB).To(spytb.Error(`expected the response to have status of 200, but it was 404`))
			})
		})
		
		t.Run("user defined status", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.WriteHeader(http.StatusTeapot)
			expect.It(t, res.Result()).To(behttp.Status(http.StatusTeapot))
		})
	})
	
	t.Run("Header matchers", func(t *testing.T) {
		t.Run("happy path multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.Header().Add("Content-Encoding", "gzip")
			res.Header().Add("Content-Type", "text/html")
			
			expect.It(t, res.Result()).To(
				behttp.Header("Content-Encoding", "gzip"),
				behttp.Header("Content-Type", "text/html"),
			)
		})
		
		t.Run("unhappy path with multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &expect.SpyTB{}
			res.Header().Add("Content-Type", "text/xml")
			
			expect.It(spyTB, res.Result()).To(
				behttp.Header("Content-Encoding", "gzip"),
				behttp.Header("Content-Type", "text/html"),
			)
			expect.It(t, spyTB).To(
				spytb.Error(`expected the response to have header "Content-Encoding" of "gzip", but it was ""`),
				spytb.Error(`expected the response to have header "Content-Type" of "text/html", but it was "text/xml"`),
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

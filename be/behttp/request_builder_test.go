package behttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be/behttp"
	"github.com/jsteenb2/expect/be/bejson"
)

func TestHTTP(t *testing.T) {
	svr := newMux()
	
	t.Run("Get", func(t *testing.T) {
		resp := behttp.Get("/").Do(svr)
		
		expect.It(t, resp).To(
			behttp.StatusOK(),
			behttp.RespBody(haveFoo(http.MethodGet)),
		)
	})
	
	t.Run("Post", func(t *testing.T) {
		resp := behttp.Post("/", nil).Do(svr)
		
		expect.It(t, resp).To(
			behttp.Status(http.StatusCreated),
			behttp.RespBody(haveFoo(http.MethodPost)),
		)
	})
	
	t.Run("Put", func(t *testing.T) {
		resp := behttp.Put("/", nil).Do(svr)
		
		expect.It(t, resp).To(
			behttp.Status(http.StatusAccepted),
			behttp.RespBody(haveFoo(http.MethodPut)),
		)
	})
	
	t.Run("Patch", func(t *testing.T) {
		resp := behttp.Patch("/", nil).Do(svr)
		
		expect.It(t, resp).To(
			behttp.Status(http.StatusPartialContent),
			behttp.RespBody(haveFoo(http.MethodPatch)),
		)
	})
	
	t.Run("Delete", func(t *testing.T) {
		resp := behttp.Delete("/").Do(svr)
		
		expect.It(t, resp).To(behttp.Status(http.StatusNoContent))
	})
	
	t.Run("Headers", func(t *testing.T) {
		resp := behttp.
			Post("/headers", strings.NewReader(`a: foo`)).
			Headers(
				"Content-Type", "text/yml",
				"Foo", "bar",
			).
			Do(svr)
		
		expect.It(t, resp).To(
			behttp.Status(http.StatusAccepted),
			behttp.ContentType("text/yml").And(behttp.Header("Foo", "bar")),
		)
	})
	
	t.Run("Queries", func(t *testing.T) {
		resp := behttp.
			Put("/queries", strings.NewReader(`a: foo`)).
			Queries(
				"ru", "bar",
				"red", "licorice",
			).
			Do(svr)
		
		expect.It(t, resp).To(
			behttp.Status(http.StatusAccepted),
			behttp.Header("ru", "bar").And(behttp.Header("red", "licorice")),
		)
	})
}

type foo struct {
	Name, Thing, Method string
}

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /headers", func(w http.ResponseWriter, req *http.Request) {
		for k := range req.Header {
			w.Header().Set(k, req.Header.Get(k))
		}
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("PUT /queries", func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		for k := range q {
			w.Header().Set(k, q.Get(k))
		}
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			writeFooResp(w, req.Method, http.StatusOK)
		case http.MethodPost:
			writeFooResp(w, req.Method, http.StatusCreated)
		case http.MethodPut:
			writeFooResp(w, req.Method, http.StatusAccepted)
		case http.MethodPatch:
			writeFooResp(w, req.Method, http.StatusPartialContent)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	return mux
}

func haveFoo(method string) expect.Matcher[io.Reader] {
	return bejson.Parsed[foo](func(got foo) expect.MatchResult {
		want := foo{
			Name:   "name",
			Thing:  "thing",
			Method: method,
		}
		return expect.MatchResult{
			Description: fmt.Sprintf("match foo %+v ", want),
			Matches:     want == got,
			But:         fmt.Sprintf("it was %+v", got),
			SubjectName: "foo response",
		}
	})
}

func writeFooResp(w http.ResponseWriter, method string, statusCode int) {
	f := foo{Name: "name", Thing: "thing", Method: method}
	writeResp(w, statusCode, f)
}

func writeResp(w http.ResponseWriter, statusCode int, v any) {
	r, err := encodeBuf(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	if _, err := io.Copy(w, r); err != nil {
		return
	}
}

func encodeBuf(v interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}

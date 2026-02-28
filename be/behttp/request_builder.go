package behttp

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// Req is a request builder.
type Req struct {
	addr    string
	method  string
	body    io.Reader
	headers []string
	queries []string
}

// Request runs creates a request for an http call. Use this if you need a
// specific/awkward request, like a GET with a request body.
func Request(method, addr string, body io.Reader) *Req {
	return &Req{
		addr:   addr,
		method: method,
		body:   body,
	}
}

// Delete creates a DELETE request.
func Delete(addr string) *Req {
	return Request(http.MethodDelete, addr, nil)
}

// Get creates a GET request.
func Get(addr string) *Req {
	return Request(http.MethodGet, addr, nil)
}

// Patch creates a PATCH request.
func Patch(addr string, body io.Reader) *Req {
	return Request(http.MethodPatch, addr, body)
}

// Post creates a POST request.
func Post(addr string, body io.Reader) *Req {
	return Request(http.MethodPost, addr, body)
}

// Put creates a PUT request.
func Put(addr string, body io.Reader) *Req {
	return Request(http.MethodPut, addr, body)
}

// Headers allows the user to set headers on the http request.
func (r *Req) Headers(k, v string, rest ...string) *Req {
	r.headers = append(r.headers, k, v)
	r.headers = append(r.headers, rest...)
	return r
}

func (r *Req) Queries(k, v string, rest ...string) *Req {
	r.queries = append(r.queries, k, v)
	r.queries = append(r.queries, rest...)
	return r
}

// Do runs the request against the provided handler.
func (r *Req) Do(handler http.Handler) *http.Response {
	req := httptest.NewRequest(r.method, r.addr, r.body)
	rec := httptest.NewRecorder()
	
	addList(req.Header, r.headers...)
	
	q := req.URL.Query()
	addList(q, r.queries...)
	req.URL.RawQuery = q.Encode()
	
	handler.ServeHTTP(rec, req)
	
	return rec.Result()
}

func addList(dst interface{ Add(k, v string) }, list ...string) {
	for i := 0; i < len(list); i += 2 {
		var v string
		if i+1 < len(list) {
			v = list[i+1]
		}
		dst.Add(list[i], v)
	}
}

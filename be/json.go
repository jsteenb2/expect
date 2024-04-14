package be

import (
	"encoding/json"
	"fmt"
	"io"
	
	"github.com/jsteenb2/expect"
)

func Parse[T any](matcher expect.Matcher[T]) expect.Matcher[io.Reader] {
	return func(rdr io.Reader) expect.MatchResult {
		var thing T
		err := json.NewDecoder(rdr).Decode(&thing)
		if err != nil {
			return expect.MatchResult{
				Description: fmt.Sprintf("be parseable into %T", thing),
				SubjectName: "JSON",
				Matches:     false,
				But:         fmt.Sprintf("it could not be parsed: %v", err),
			}
		}
		return matcher(thing)
	}
}

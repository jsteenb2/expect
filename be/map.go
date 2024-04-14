package be

import (
	"fmt"
	
	"github.com/jsteenb2/expect"
)

// HaveKey checks if a map has a key with a specific value.
func HaveKey[K comparable, V any](key K, valueMatcher expect.Matcher[V]) expect.Matcher[map[K]V] {
	return func(m map[K]V) expect.MatchResult {
		value, exists := m[key]
		
		if !exists {
			return missingKeyResult(key)
		}
		
		result := valueMatcher(value)
		result.Description = fmt.Sprintf("have key %v with value %v", key, result.Description)
		result.SubjectName = fmt.Sprintf("%+v", m)
		return result
	}
}

// WithAnyValue lets you match any value, useful if you're just looking for the presence of a key.
func WithAnyValue[T any]() expect.Matcher[T] {
	return func(T) expect.MatchResult {
		return expect.MatchResult{
			Matches: true,
		}
	}
}

func missingKeyResult[K any](key K) expect.MatchResult {
	return expect.MatchResult{
		Description: fmt.Sprintf("have key %v", key),
		Matches:     false,
		But:         "it did not",
	}
}

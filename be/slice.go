package be

import (
	"fmt"
	"slices"
	
	"github.com/jsteenb2/expect"
)

var passingResult = expect.MatchResult{
	Matches: true,
}

// HaveSize checks if an array's size meets a matcher's criteria.
func HaveSize[T any](matcher expect.Matcher[int]) expect.Matcher[[]T] {
	return func(items []T) expect.MatchResult {
		result := matcher(len(items))
		result.Description = "have a size " + result.Description
		return result
	}
}

// ContainItem checks if an array contains an item that meets a matcher's criteria.
func ContainItem[T any](m expect.Matcher[T]) expect.Matcher[[]T] {
	return func(items []T) expect.MatchResult {
		var exampleFailure expect.MatchResult
		
		for _, item := range items {
			result := m(item)
			if result.Matches {
				return expect.MatchResult{
					Description: "contain an item",
					Matches:     true,
				}
			} else {
				exampleFailure = result
			}
		}
		
		exampleFailure.But = "it did not"
		exampleFailure.Description = "contain an item " + exampleFailure.Description
		exampleFailure.SubjectName = fmt.Sprintf("%+v", items)
		
		return exampleFailure
	}
}

// EveryItem checks if every item in an array meets a matcher's criteria.
func EveryItem[T any](m expect.Matcher[T]) expect.Matcher[[]T] {
	return func(items []T) expect.MatchResult {
		for _, item := range items {
			if result := m(item); !result.Matches {
				return everyItemFailure(result)
			}
		}
		
		return passingResult
	}
}

// ShallowEquals checks if two slices are equal, only works with slices of comparable types.
func ShallowEquals[T comparable](other []T) expect.Matcher[[]T] {
	return func(ts []T) expect.MatchResult {
		return expect.MatchResult{
			Matches:     slices.Equal(ts, other),
			Description: fmt.Sprintf("be equal to %v", other),
		}
	}
}

func everyItemFailure(result expect.MatchResult) expect.MatchResult {
	return expect.MatchResult{
		Description: "have every item " + result.Description,
		Matches:     false,
		But:         result.But,
	}
}

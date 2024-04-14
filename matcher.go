package expect

type Matcher[T any] func(T) MatchResult

// Or combines matchers with a boolean OR.
func (m Matcher[T]) Or(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)
		
		for _, matcher := range matchers {
			result.Description += " or " + matcher(got).Description
		}
		
		if result.Matches {
			return result
		}
		
		for _, matcher := range matchers {
			if r := matcher(got); r.Matches {
				result.Matches = true
				return result
			}
		}
		
		return result
	}
}

// And combines matchers with a boolean AND.
func (m Matcher[T]) And(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)
		
		for _, matcher := range matchers {
			r := matcher(got)
			result = result.Combine(r)
		}
		
		return result
	}
}

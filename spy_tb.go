package expect

import (
	"fmt"
	"strconv"
	"strings"
)

// SpyTB is a test helper that records calls to Error. This lets us create tests and examples to demonstrate what happens when a test fails.
type SpyTB struct {
	ErrorCalls []string
}

// Result will print the test result.
func (s *SpyTB) Result() string {
	if len(s.ErrorCalls) == 0 {
		return "Test passed"
	}
	return fmt.Sprintf("Test failed: %v", s.ErrorCalls)
}

func (s *SpyTB) String() string {
	return "Spy TB"
}

func (s *SpyTB) Helper() {
}

func (s *SpyTB) Error(args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprint(args...))
}

func (s *SpyTB) Errorf(format string, args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprintf(format, args...))
}

func (s *SpyTB) Fatalf(format string, args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprintf(format, args...))
}

func (s *SpyTB) Reset() {
	s.ErrorCalls = nil
}

func (s *SpyTB) Format(f fmt.State, verb rune) {
	switch verb {
	case 'q', 's':
		if len(s.ErrorCalls) == 0 {
			fmt.Fprint(f, "Test passed")
			return
		}
		for i := range s.ErrorCalls {
			s.ErrorCalls[i] = stripStackTrace(s.ErrorCalls[i])
			if verb == 'q' {
				s.ErrorCalls[i] = strconv.Quote(s.ErrorCalls[i])
			}
		}
		fmt.Fprintf(f, "Test failed: %v", s.ErrorCalls)
	case 'v':
		if f.Flag('#') {
			fmt.Fprintf(f, "%#v", s.ErrorCalls)
			return
		}
		if len(s.ErrorCalls) == 0 {
			fmt.Fprint(f, "Test passed")
			return
		}
		fmt.Fprintf(f, "Test failed: %v", s.ErrorCalls)
	default:
		fmt.Fprint(f, s.Result())
	}
}

func stripStackTrace(err string) string {
	idx := strings.Index(err, "Error Trace:")
	if idx == -1 {
		return err
	}
	return strings.TrimSpace(err[:idx])
}

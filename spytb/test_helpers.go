package spytb

import (
	"github.com/jsteenb2/expect"
)

func VerifyFailingMatcher[T any](t expect.TB, subject T, matcher expect.Matcher[T], expectedError string) {
	t.Helper()
	spyTB := &expect.SpyTB{}
	expect.It(spyTB, subject).To(matcher)
	expect.It(t, spyTB).To(Error(expectedError))
}

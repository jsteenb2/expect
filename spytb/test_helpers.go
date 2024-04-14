package spytb

import (
	"github.com/jsteenb2/expect"
)

func VerifyFailingMatcher[T any](t expect.TB, subject T, matcher expect.Matcher[T], expectedError string) {
	t.Helper()
	spyTB := &expect.SpyTB{}
	expect.Expect(spyTB, subject).To(matcher)
	expect.Expect(t, spyTB).To(HaveError(expectedError))
}

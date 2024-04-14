package spytb

import (
	"testing"
	
	"github.com/jsteenb2/expect"
)

func TestSpyTB(t *testing.T) {
	t.Run("correctly had errors", func(t *testing.T) {
		spyTB := &expect.SpyTB{}
		subject := &expect.SpyTB{ErrorCalls: []string{"oh no"}}
		
		expect.Expect(spyTB, subject).To(HaveError("oopsie"))
		expect.Expect(t, spyTB).To(HaveError(`expected Spy TB to have error "oopsie", but has [oh no]`))
	})
	
	t.Run("complains if it has errors when none expected", func(t *testing.T) {
		spyTB := &expect.SpyTB{}
		subject := &expect.SpyTB{ErrorCalls: []string{"oh no"}}
		
		expect.Expect(spyTB, subject).To(HaveNoErrors)
		expect.Expect(t, spyTB).To(HaveError(`expected Spy TB to have no errors, but it had errors [oh no]`))
	})
}
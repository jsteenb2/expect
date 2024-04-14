package spytb

import (
	"testing"
	
	"github.com/jsteenb2/expect"
)

func TestSpyTB(t *testing.T) {
	t.Run("correctly had errors", func(t *testing.T) {
		spyTB := &expect.SpyTB{}
		subject := &expect.SpyTB{ErrorCalls: []string{"oh no"}}
		
		expect.It(spyTB, subject).To(Error("oopsie"))
		expect.It(t, spyTB).To(Error(`expected Spy TB to have error "oopsie", but has [oh no]`))
	})
	
	t.Run("complains if it has errors when none expected", func(t *testing.T) {
		spyTB := &expect.SpyTB{}
		subject := &expect.SpyTB{ErrorCalls: []string{"oh no"}}
		
		expect.It(spyTB, subject).To(NoErrors)
		expect.It(t, spyTB).To(Error(`expected Spy TB to have no errors, but it had errors [oh no]`))
	})
}

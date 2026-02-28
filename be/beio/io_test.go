package beio_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/be/beio"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleContainingByte() {
	t := &expect.SpyTB{}
	
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	
	expect.It[io.Reader](t, buf).To(beio.HaveData(
		beio.ContainingByte([]byte("hello")).And(beio.ContainingByte([]byte("world"))),
	))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleContainingByte_fail() {
	t := &expect.SpyTB{}
	
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	
	expect.It[io.Reader](t, buf).To(beio.HaveData(
		beio.ContainingByte([]byte("goodbye")),
	))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the reader to contain "goodbye", but it didn't have "goodbye"]
}

func ExampleContainingString() {
	t := &expect.SpyTB{}
	
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	
	expect.It[io.Reader](t, buf).To(beio.HaveData(
		beio.ContainingString("world"),
	))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleContainingString_fail() {
	t := &expect.SpyTB{}
	
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	
	expect.It[io.Reader](t, buf).To(beio.HaveData(
		beio.ContainingString("goodbye"),
	))
	fmt.Println(t.Result())
	// Output: Test failed: [expected the reader to contain "goodbye", but it was "helloworld"]
}

func ExampleString() {
	t := &expect.SpyTB{}
	
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	expect.It[io.Reader](t, buf).To(beio.String(be.Eq("helloworld")))
	fmt.Println(t.Result())
	// Output: Test passed
}

func ExampleString_fail() {
	t := &expect.SpyTB{}
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	expect.It[io.Reader](t, buf).To(beio.String(be.Eq("Poo")))
	fmt.Println(t.Result())
	// Output: Test failed: [expected "helloworld" to be equal to "Poo", but it was "helloworld"]
}

func TestIOMatchers(t *testing.T) {
	t.Run("passing", func(t *testing.T) {
		buf := &bytes.Buffer{}
		buf.WriteString("hello")
		buf.WriteString("world")
		
		expect.It[io.Reader](t, buf).To(beio.HaveData(
			beio.ContainingByte([]byte("hello")).And(beio.ContainingByte([]byte("world"))),
		))
		
		buf.WriteString("goodbye")
		expect.It[io.Reader](t, buf).To(beio.HaveData(
			beio.ContainingString("goodbye"),
		))
	})
	
	t.Run("failing", func(t *testing.T) {
		buf := &bytes.Buffer{}
		buf.WriteString("hello")
		buf.WriteString("world")
		
		spytb.VerifyFailingMatcher[io.Reader](
			t,
			buf,
			beio.HaveData(beio.ContainingByte([]byte("goodbye"))),
			`expected the reader to contain "goodbye", but it didn't have "goodbye"`,
		)
		
		spytb.VerifyFailingMatcher[io.Reader](
			t,
			buf,
			beio.HaveData(beio.ContainingString("goodbye")),
			`expected the reader to contain "goodbye", but it was ""`,
		)
	})
}

package befs_test

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"
	
	"github.com/jsteenb2/expect"
	"github.com/jsteenb2/expect/be"
	"github.com/jsteenb2/expect/be/befs"
	"github.com/jsteenb2/expect/be/beio"
	"github.com/jsteenb2/expect/spytb"
)

func ExampleFileNamed_fail() {
	t := &expect.SpyTB{}
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
	}
	
	expect.It[fs.FS](t, stubFS).To(befs.FileNamed("someFile.txt", beio.String(be.Substring("Pluto"))))

	fmt.Printf("%s\n", t)
	// Output: Test failed: [expected file called "someFile.txt" to contain "Pluto", but while the file existed, the contents did not match]
}

func ExampleFileNamed() {
	t := &expect.SpyTB{}
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
	}
	
	expect.It[fs.FS](t, stubFS).To(befs.FileNamed("someFile.txt"))

	fmt.Printf("%s\n", t)
	// Output: Test passed
}

func ExampleDir() {
	t := &expect.SpyTB{}
	stubFS := fstest.MapFS{
		"someDir": {
			Mode: fs.ModeDir,
		},
	}
	
	expect.It[fs.FS](t, stubFS).To(befs.Dir("someDir"))

	fmt.Printf("%s\n", t)
	// Output: Test passed
}

func ExampleDir_fail() {
	t := &expect.SpyTB{}
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
	}
	
	expect.It[fs.FS](t, stubFS).To(befs.Dir("someFile.txt"))

	fmt.Printf("%s\n", t)
	// Output: Test failed: [expected file system to have directory called "someFile.txt", but it was not a directory]
}

func TestFSMatching(t *testing.T) {
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
		"someDir": {
			Mode: fs.ModeDir,
		},
		"nested/someFile.txt": {
			Data: []byte("hello world"),
		},
	}
	
	t.Run("HasDir", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It[fs.FS](t, stubFS).To(befs.Dir("someDir"))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				befs.Dir("someFile.txt"),
				`expected file system to have directory called "someFile.txt", but it was not a directory`,
			)
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				befs.Dir("non-existent-file"),
				`expected file system to have directory called "non-existent-file", but it did not`,
			)
			t.Run("failing filesystem", func(t *testing.T) {
				failingFS := FailToReadFS{Error: fmt.Errorf("could not read file")}
				spytb.VerifyFailingMatcher[fs.FS](
					t,
					failingFS,
					befs.Dir("someDir"),
					`expected file system to have directory called "someDir", but it could not be read`,
				)
			})
		})
	})
	
	t.Run("FileContains", func(t *testing.T) {
		t.Run("file existence check", func(t *testing.T) {
			t.Run("passing", func(t *testing.T) {
				expect.It[fs.FS](t, stubFS).To(befs.FileNamed("someFile.txt"))
				expect.It[fs.FS](t, stubFS).To(befs.FileNamed("nested/someFile.txt"))
			})
			t.Run("failing", func(t *testing.T) {
				spytb.VerifyFailingMatcher[fs.FS](
					t,
					stubFS,
					befs.FileNamed("non-existent-file"),
					`expected file system to have file called "non-existent-file", but it did not`,
				)
			})
		})
	})
	
	t.Run("FileContains with contents", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			expect.It[fs.FS](t, stubFS).To(befs.FileNamed("someFile.txt", beio.String(be.Substring("world"))))
		})
		
		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				befs.FileNamed("someFile.txt", beio.String(be.Substring("goodbye"))),
				`expected file called "someFile.txt" to contain "goodbye"`,
			)
			
			t.Run("failing filesystem", func(t *testing.T) {
				failingFS := FailToReadFS{Error: fmt.Errorf("could not read file")}
				spytb.VerifyFailingMatcher[fs.FS](
					t,
					failingFS,
					befs.FileNamed("anotherFile.txt", beio.String(be.Substring("BLAH"))),
					`expected file called "anotherFile.txt" to have data in io.Reader, but it could not be read`,
				)
			})
		})
	})
}

type FailToReadFS struct {
	Error error
}

func (f FailToReadFS) Open(name string) (fs.File, error) {
	return FailingFile(f), nil
}

type FailingFile struct {
	Error error
}

func (f FailingFile) Stat() (fs.FileInfo, error) {
	return nil, f.Error
}

func (f FailingFile) Read(bytes []byte) (int, error) {
	return 0, f.Error
}

func (f FailingFile) Close() error {
	return f.Error
}

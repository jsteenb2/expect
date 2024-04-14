package be

import (
	"fmt"
	"io"
	"io/fs"
	
	"github.com/jsteenb2/expect"
)

const subjectName = "file system"

// HaveFileCalled checks if a file exists in the file system, and can run additional matchers on its contents.
func HaveFileCalled(name string, contentMatcher ...expect.Matcher[string]) expect.Matcher[fs.FS] {
	return func(fileSystem fs.FS) expect.MatchResult {
		file, err := fileSystem.Open(name)
		
		if err != nil {
			return expect.MatchResult{
				Description: "have file called " + name,
				Matches:     false,
				But:         "it did not",
				SubjectName: subjectName,
			}
		}
		
		defer file.Close()
		
		if len(contentMatcher) > 0 {
			all, err := io.ReadAll(file)
			if err != nil {
				return expect.MatchResult{
					Description: "have file called " + name,
					Matches:     false,
					But:         "it could not be read",
					SubjectName: subjectName,
				}
			}
			contents := string(all)
			for _, matcher := range contentMatcher {
				result := matcher(contents)
				result.SubjectName = "file called " + name
				if !result.Matches {
					return result
				}
			}
		}
		
		return expect.MatchResult{
			Description: "have file called " + name,
			Matches:     true,
			SubjectName: subjectName,
		}
	}
}

// HaveDir checks if a directory exists in the file system.
func HaveDir(name string) expect.Matcher[fs.FS] {
	return func(fileSystem fs.FS) expect.MatchResult {
		f, err := fileSystem.Open(name)
		
		result := expect.MatchResult{
			Description: fmt.Sprintf("have directory called %q", name),
			SubjectName: subjectName,
			Matches:     true,
		}
		
		if err != nil {
			result.Matches = false
			result.But = "it did not"
			return result
		}
		
		stat, err := f.Stat()
		if err != nil {
			result.Matches = false
			result.But = "it could not be read"
			return result
		}
		
		if !stat.IsDir() {
			result.Matches = false
			result.But = "it was not a directory"
			return result
		}
		
		return result
	}
}

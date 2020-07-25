package gitignore

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type assert struct {
	patterns []string
	file     file
	expect   bool
}

type file struct {
	path  string
	isDir bool
}

func TestFromFile(t *testing.T) {
	fp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(fp.Name())
	lines := []string{"a.txt", "b.txt"}
	_, _ = fp.WriteString(strings.Join(lines, "\n"))
	_, err = FromFile(fp.Name())
	if err != nil {
		t.Error(err)
	}
	_, err = FromFile(fp.Name(), ".")
	if err != nil {
		t.Error(err)
	}
}

func TestDummyMatcher(t *testing.T) {
	if AllowedMatcher.Match("", false) {
		t.Error("AllowedMatcher not expected")
	}
	if !IgnoredMatcher.Match("", false) {
		t.Error("IgnoredMatcher not expected")
	}
}

func TestCombine(t *testing.T) {
	asserts := map[file]bool{
		file{".DS_Store", false}: true,
		file{"foo.txt", false}:   false,
	}
	combine := Combine(
		FromLines(".", []string{".DS_Store"}),
		FromLines(".", []string{"!/foo.txt"}),
	)
	for file, expect := range asserts {
		result := combine.Match(file.path, file.isDir)
		if result != expect {
			t.Errorf("Match should return %t, got %t", expect, result)
		}
	}
}

func TestMatch(t *testing.T) {
	asserts := []assert{
		{[]string{"a.txt"}, file{"a.txt", false}, true},
		{[]string{"*.txt"}, file{"a.txt", false}, true},
		{[]string{"dir/a.txt"}, file{"dir/a.txt", false}, true},
		{[]string{"dir/*.txt"}, file{"dir/a.txt", false}, true},
		{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		{[]string{"dir3/a.txt"}, file{"dir1/dir2/dir3/a.txt", false}, true},
		{[]string{"a.txt"}, file{"dir/a.txt", false}, true},
		{[]string{"*.txt"}, file{"dir/a.txt", false}, true},
		{[]string{"a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		{[]string{"dir"}, file{"dir", true}, true},
		{[]string{"dir/"}, file{"dir", true}, true},
		{[]string{"dir/"}, file{"dir", false}, false},
		{[]string{"dir1/dir2/"}, file{"dir1/dir2", true}, true},
		{[]string{"/a.txt"}, file{"a.txt", false}, true},
		{[]string{"/dir/a.txt"}, file{"dir/a.txt", false}, true},
		{[]string{"/dir1/a.txt"}, file{"dir/dir1/a.txt", false}, false},
		{[]string{"/a.txt"}, file{"dir/a.txt", false}, false},
		{[]string{"a.txt", "b.txt"}, file{"dir/b.txt", false}, true},
		{[]string{"*.txt", "!b.txt"}, file{"dir/b.txt", false}, false},
		{[]string{"dir/*.txt", "!dir/b.txt"}, file{"dir/b.txt", false}, false},
		{[]string{"dir/*.txt", "!/b.txt"}, file{"dir/b.txt", false}, true},
		{[]string{`\#a.txt`}, file{"#a.txt", false}, true},
	}

	for _, assert := range asserts {
		gi := FromLines(".", assert.patterns)
		result := gi.Match(assert.file.path, assert.file.isDir)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}
	}
}

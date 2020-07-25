package gitignore

import (
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
		gi := NewGitIgnoreFromLines(".", assert.patterns)
		result := gi.Match(assert.file.path, assert.file.isDir)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}
	}
}

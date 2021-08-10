package gitignore

import (
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

func TestMatch(t *testing.T) {
	asserts := []assert{
		assert{[]string{"a.txt"}, file{"a.txt", false}, true},
		assert{[]string{"*.txt"}, file{"a.txt", false}, true},
		assert{[]string{"dir/a.txt"}, file{"dir/a.txt", false}, true},
		assert{[]string{"dir/*.txt"}, file{"dir/a.txt", false}, true},
		assert{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		assert{[]string{"dir3/a.txt"}, file{"dir1/dir2/dir3/a.txt", false}, true},
		assert{[]string{"a.txt"}, file{"dir/a.txt", false}, true},
		assert{[]string{"*.txt"}, file{"dir/a.txt", false}, true},
		assert{[]string{"a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		assert{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false}, true},
		assert{[]string{"dir"}, file{"dir", true}, true},
		assert{[]string{"dir/"}, file{"dir", true}, true},
		assert{[]string{"dir/"}, file{"dir", false}, false},
		assert{[]string{"dir1/dir2/"}, file{"dir1/dir2", true}, true},
		assert{[]string{"/a.txt"}, file{"a.txt", false}, true},
		assert{[]string{"/dir/a.txt"}, file{"dir/a.txt", false}, true},
		assert{[]string{"/dir1/a.txt"}, file{"dir/dir1/a.txt", false}, false},
		assert{[]string{"/a.txt"}, file{"dir/a.txt", false}, false},
		assert{[]string{"a.txt", "b.txt"}, file{"dir/b.txt", false}, true},
		assert{[]string{"*.txt", "!b.txt"}, file{"dir/b.txt", false}, false},
		assert{[]string{"dir/*.txt", "!dir/b.txt"}, file{"dir/b.txt", false}, false},
		assert{[]string{"dir/*.txt", "!/b.txt"}, file{"dir/b.txt", false}, true},
		assert{[]string{`\#a.txt`}, file{"#a.txt", false}, true},
		assert{[]string{`**/__test__/**`}, file{"dir1/dir2/__test__/dir3/b.txt", false}, true},
	}

	for _, assert := range asserts {
		gi := NewGitIgnoreFromReader(".", strings.NewReader(strings.Join(assert.patterns, "\n")))
		result := gi.Match(assert.file.path, assert.file.isDir)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}
	}
}

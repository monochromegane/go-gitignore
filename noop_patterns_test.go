package gitignore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoopPatterns(t *testing.T) {
	base, err := ioutil.TempDir("", "gitignore-noop-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(base)

	for _, tc := range []struct {
		name     string
		patterns string
		logs     bool
	}{
		{"empty", "", false},
		{"root", "/", false},
		{"negation", "!", false},
		{"negated-root", "!/", false},
		{"multiple", "/\n!\n!/\n//\n///\n", false},
		{"root-before-rules", "/\n*.log\n!keep.log\n", true},
		{"root-after-rules", "*.log\n!keep.log\n/", true},
		{"negation-between-rules", "*.log\n!\n!/\n!keep.log\n", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, fromFile := range []bool{false, true} {
				name := "reader"
				if fromFile {
					name = "file"
				}
				t.Run(name, func(t *testing.T) {
					var matcher IgnoreMatcher
					if fromFile {
						path := filepath.Join(base, ".gitignore")
						if err := ioutil.WriteFile(path, []byte(tc.patterns), 0600); err != nil {
							t.Fatal(err)
						}
						var err error
						matcher, err = NewGitIgnore(path)
						if err != nil {
							t.Fatal(err)
						}
					} else {
						matcher = NewGitIgnoreFromReader(base, strings.NewReader(tc.patterns))
					}
					for _, entry := range []struct {
						path  string
						isDir bool
						log   bool
					}{
						{"a.txt", false, false},
						{"a.log", false, true},
						{"keep.log", false, false},
						{"dir", true, false},
						{"dir/a.log", false, true},
						{"dir/keep.log", false, false},
					} {
						want := tc.logs && entry.log
						if got := matcher.Match(filepath.Join(base, entry.path), entry.isDir); got != want {
							t.Errorf("Match(%q, %t) = %t, want %t", entry.path, entry.isDir, got, want)
						}
					}
				})
			}
		})
	}
}

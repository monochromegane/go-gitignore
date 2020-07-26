package gitignore

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type IgnoreMatcher interface {
	Match(path string, isDir bool) bool
}

type IgnoreMatcherFunc func(path string, isDir bool) bool

func (m IgnoreMatcherFunc) Match(path string, isDir bool) bool {
	return m(path, isDir)
}

var AllowedMatcher = IgnoreMatcherFunc(func(string, bool) bool { return false })
var IgnoredMatcher = IgnoreMatcherFunc(func(string, bool) bool { return true })

type gitIgnore struct {
	ignorePatterns scanStrategy
	acceptPatterns scanStrategy
	path           string
}

func (g *gitIgnore) Match(path string, isDir bool) bool {
	path, err := filepath.Rel(g.path, path)
	if err != nil {
		return false
	}
	path = filepath.ToSlash(path)
	if g.acceptPatterns.match(path, isDir) {
		return false
	}
	return g.ignorePatterns.match(path, isDir)
}

func NewFromFile(pattern string, base ...string) (matcher IgnoreMatcher, err error) {
	path := filepath.Dir(pattern)
	if base != nil {
		path = filepath.Join(base...)
	}
	file, err := os.Open(pattern)
	if err == nil {
		defer file.Close()
		matcher = NewFromReader(path, file)
	}
	return
}

func NewFromReader(path string, r io.Reader) IgnoreMatcher {
	scanner := bufio.NewScanner(r)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return NewFromLines(path, lines)
}

func NewFromLines(path string, lines []string) IgnoreMatcher {
	matcher := &gitIgnore{
		ignorePatterns: newIndexScanPatterns(),
		acceptPatterns: newIndexScanPatterns(),
		path:           path,
	}
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if len(line) > 2 && line[:2] == `\#` {
			line = line[1:]
		}
		if line[0] == '!' {
			matcher.acceptPatterns.add(line[1:])
		} else {
			matcher.ignorePatterns.add(line)
		}
	}
	return matcher
}

func Combine(matchers ...IgnoreMatcher) IgnoreMatcher {
	return IgnoreMatcherFunc(func(path string, isDir bool) bool {
		for _, matcher := range matchers {
			if matcher.Match(path, isDir) {
				return true
			}
		}
		return false
	})
}

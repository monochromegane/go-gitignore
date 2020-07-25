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

type IgnoreMatcherFn func(path string, isDir bool) bool

func (m IgnoreMatcherFn) Match(path string, isDir bool) bool {
	return m(path, isDir)
}

var AllowedMatcher = IgnoreMatcherFn(func(string, bool) bool { return false })
var IgnoredMatcher = IgnoreMatcherFn(func(string, bool) bool { return true })

type gitIgnore struct {
	ignorePatterns scanStrategy
	acceptPatterns scanStrategy
	path           string
}

func FromFile(gitignore string, base ...string) (IgnoreMatcher, error) {
	var path string
	if len(base) > 0 {
		path = base[0]
	} else {
		path = filepath.Dir(gitignore)
	}

	file, err := os.Open(gitignore)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return FromReader(path, file), nil
}

func FromReader(path string, r io.Reader) IgnoreMatcher {
	scanner := bufio.NewScanner(r)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		lines = append(lines, line)
	}
	return FromLines(path, lines)
}

func FromLines(path string, lines []string) IgnoreMatcher {
	g := gitIgnore{
		ignorePatterns: newIndexScanPatterns(),
		acceptPatterns: newIndexScanPatterns(),
		path:           path,
	}
	for _, line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, `\#`) {
			line = strings.TrimPrefix(line, `\`)
		}
		if strings.HasPrefix(line, "!") {
			g.acceptPatterns.add(strings.TrimPrefix(line, "!"))
		} else {
			g.ignorePatterns.add(line)
		}
	}
	return g
}

func Combine(matchers ...IgnoreMatcher) IgnoreMatcher {
	return IgnoreMatcherFn(func(path string, isDir bool) bool {
		for _, matcher := range matchers {
			if matcher.Match(path, isDir) {
				return true
			}
		}
		return false
	})
}

func (g gitIgnore) Match(path string, isDir bool) bool {
	relativePath, err := filepath.Rel(g.path, path)
	if err != nil {
		return false
	}
	relativePath = filepath.ToSlash(relativePath)

	if g.acceptPatterns.match(relativePath, isDir) {
		return false
	}
	return g.ignorePatterns.match(relativePath, isDir)
}

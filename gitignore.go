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

type gitIgnore struct {
	ignorePatterns patterns
	acceptPatterns patterns
	path           string
}

func NewGitIgnore(gitignore string) (IgnoreMatcher, error) {
	path := filepath.Dir(gitignore)
	file, err := os.Open(gitignore)
	if err != nil {
		return nil, err
	}

	return newGitIgnore(path, file), nil
}

func newGitIgnore(path string, r io.Reader) gitIgnore {
	g := gitIgnore{path: path}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			g.acceptPatterns = append(g.acceptPatterns,
				newPattern(strings.TrimPrefix(line, "!"), g.path))
		} else {
			g.ignorePatterns = append(g.ignorePatterns, newPattern(line, g.path))
		}
	}
	return g
}

func (g gitIgnore) Match(path string, isDir bool) bool {
	if match := g.acceptPatterns.match(path, isDir); match {
		return false
	}
	return g.ignorePatterns.match(path, isDir)
}

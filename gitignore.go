package gitignore

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type gitIgnore struct {
	ignorePatterns patterns
	acceptPatterns patterns
	path           string
}

func NewGitIgnore(gitignore string) gitIgnore {
	path := filepath.Dir(gitignore)
	file, _ := os.Open(gitignore)
	return newGitIgnore(path, file)
}

func newGitIgnore(path string, r io.Reader) gitIgnore {

	g := gitIgnore{path: path}

	if r == nil {
		return g
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			g.acceptPatterns = append(g.acceptPatterns,
				pattern{strings.TrimPrefix(line, "!"), g.path})
		} else {
			g.ignorePatterns = append(g.ignorePatterns, pattern{line, g.path})
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

package gitignore

import (
	"path/filepath"
	"strings"
)

var Separator = string(filepath.Separator)

type pattern struct {
	path string
	base string
}

func (p pattern) match(path string, isDir bool) bool {

	if p.hasDirSuffix() && !isDir {
		return false
	}

	pattern := p.trimedPattern()

	var match bool
	if p.hasRootPrefix() {
		// absolute pattern
		match, _ = filepath.Match(filepath.Join(p.base, p.path), path)
	} else {
		// relative pattern
		match, _ = filepath.Match(pattern, p.equalizeDepth(path))
	}
	return match
}

func (p pattern) equalizeDepth(path string) string {
	trimedPath := strings.TrimPrefix(path, p.base)
	patternDepth := strings.Count(p.path, "/")
	pathDepth := strings.Count(trimedPath, Separator)
	start := 0
	if diff := pathDepth - patternDepth; diff > 0 {
		start = diff
	}
	return filepath.Join(strings.Split(trimedPath, Separator)[start:]...)
}

func (p pattern) prefix() string {
	return string(p.path[0])
}

func (p pattern) suffix() string {
	return string(p.path[len(p.path)-1])
}

func (p pattern) hasRootPrefix() bool {
	return p.prefix() == "/"
}

func (p pattern) hasDirSuffix() bool {
	return p.suffix() == "/"
}

func (p pattern) trimedPattern() string {
	return strings.Trim(p.path, "/")
}

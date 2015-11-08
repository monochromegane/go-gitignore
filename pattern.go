package gitignore

import (
	"path/filepath"
	"strings"
)

var Separator = string(filepath.Separator)

type pattern struct {
	path          string
	base          string
	matchingPath  string
	hasRootPrefix bool
	hasDirSuffix  bool
	pathDepth     int
}

func newPattern(path, base string) pattern {
	hasRootPrefix := path[0] == '/'
	hasDirSuffix := path[len(path)-1] == '/'

	var matchingPath string
	var pathDepth int
	if hasRootPrefix {
		matchingPath = filepath.Join(base, path)
	} else {
		matchingPath = strings.Trim(path, "/")
		pathDepth = strings.Count(path, "/")
	}

	return pattern{
		path:          path,
		base:          base,
		matchingPath:  matchingPath,
		hasRootPrefix: hasRootPrefix,
		hasDirSuffix:  hasDirSuffix,
		pathDepth:     pathDepth,
	}
}

func (p pattern) match(path string, isDir bool) bool {

	if p.hasDirSuffix && !isDir {
		return false
	}

	var match bool
	if p.hasRootPrefix {
		//absolute pattern
		match, _ = filepath.Match(p.matchingPath, path)
	} else {
		// relative pattern
		match, _ = filepath.Match(p.matchingPath, p.equalizeDepth(path))
	}
	return match
}

func (p pattern) equalizeDepth(path string) string {
	trimedPath := strings.TrimPrefix(path, p.base)
	pathDepth := strings.Count(trimedPath, Separator)
	start := 0
	if diff := pathDepth - p.pathDepth; diff > 0 {
		start = diff
	}
	return filepath.Join(strings.Split(trimedPath, Separator)[start:]...)
}

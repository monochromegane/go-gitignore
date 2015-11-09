package gitignore

import (
	"path/filepath"
	"strings"
)

var Separator = string(filepath.Separator)

type pattern struct {
	path          string
	base          string
	hasRootPrefix bool
	hasDirSuffix  bool
	pathDepth     int
	matcher       pathMatcher
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

	var matcher pathMatcher
	if hasMeta(path) {
		matcher = filepathMatcher{path: matchingPath}
	} else {
		matcher = simpleMatcher{path: matchingPath}
	}

	return pattern{
		base:          base,
		hasRootPrefix: hasRootPrefix,
		hasDirSuffix:  hasDirSuffix,
		pathDepth:     pathDepth,
		matcher:       matcher,
	}
}

func (p pattern) match(path string, isDir bool) bool {

	if p.hasDirSuffix && !isDir {
		return false
	}

	var targetPath string
	if p.hasRootPrefix {
		//absolute pattern
		targetPath = path
	} else {
		// relative pattern
		targetPath = p.equalizeDepth(path)
	}
	return p.matcher.match(targetPath)
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

func hasMeta(path string) bool {
	return strings.IndexAny(path, "*?[") >= 0
}

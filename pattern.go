package gitignore

import (
	"os"
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
	return cutLastN(trimedPath, p.pathDepth+1)
}

func cutLastN(path string, n int) string {
	vol := filepath.VolumeName(path)
	i := len(path) - 1

	var count int
	for i >= len(vol) {
		if os.IsPathSeparator(path[i]) {
			count++
			if count >= n {
				break
			}
		}
		i--
	}
	return path[i+1:]
}

func hasMeta(path string) bool {
	return strings.IndexAny(path, "*?[") >= 0
}

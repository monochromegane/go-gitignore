package gitignore

import "strings"

type indexedPatterns struct {
	absolute depthPatternHolder
	relative depthPatternHolder
}

func newIndexedPatterns(base string) indexedPatterns {
	return indexedPatterns{
		absolute: newDepthPatternHolder(asc, base),
		relative: newDepthPatternHolder(desc, base),
	}
}

func (ps *indexedPatterns) add(pattern string) {
	if strings.HasPrefix(pattern, "/") {
		ps.absolute.add(pattern)
	} else {
		ps.relative.add(pattern)
	}
}

func (ps indexedPatterns) match(path string, isDir bool) bool {
	if ps.absolute.match(path, isDir) {
		return true
	}
	return ps.relative.match(path, isDir)
}

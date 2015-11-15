package gitignore

import "strings"

const initials = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type initialPatternHolder struct {
	patterns      initialPatterns
	otherPatterns patterns
	base          string
}

func newInitialPatternHolder(base string) initialPatternHolder {
	return initialPatternHolder{
		patterns:      initialPatterns{m: map[byte]patterns{}},
		otherPatterns: patterns{},
		base:          base,
	}
}

func (h *initialPatternHolder) add(pattern string) {
	if strings.IndexAny(pattern[0:1], initials) != -1 {
		h.patterns.set(pattern[0], newPatternForEqualizedPath(pattern, h.base))
	} else {
		h.otherPatterns = append(h.otherPatterns, newPatternForEqualizedPath(pattern, h.base))
	}
}

func (h initialPatternHolder) match(path string, isDir bool) bool {
	if patterns, ok := h.patterns.get(path[0]); ok {
		if patterns.match(path, isDir) {
			return true
		}
	}
	return h.otherPatterns.match(path, isDir)
}

type initialPatterns struct {
	m map[byte]patterns
}

func (p *initialPatterns) set(initial byte, pattern pattern) {
	if ps, ok := p.m[initial]; ok {
		p.m[initial] = append(ps, pattern)
	} else {
		p.m[initial] = patterns{pattern}
	}
}

func (p initialPatterns) get(initial byte) (patterns, bool) {
	patterns, ok := p.m[initial]
	return patterns, ok
}

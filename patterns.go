package gitignore

type patterns []pattern

func (ps patterns) match(path string, isDir bool) bool {
	for _, p := range ps {
		if match := p.match(path, isDir); match {
			return true
		}
	}
	return false
}

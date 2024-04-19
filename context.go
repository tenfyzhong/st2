package st2

// Context struct contains the context running
type Context struct {
	Src    string
	Dst    string
	Root   string
	Prefix string
	Suffix string
}

func NewContext(src, dst, root, prefix, suffix string) Context {
	return Context{
		Src:    src,
		Dst:    dst,
		Root:   normalizeToken(root, ""),
		Prefix: normalizeToken(prefix, ""),
		Suffix: normalizeToken(suffix, ""),
	}
}

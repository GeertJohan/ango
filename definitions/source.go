package definitions

// Source contains information about where a given definition was declared.
type Source struct {
	// Filename   string // no support for filename, parser is singlefile now anyway

	Linenumber int
}

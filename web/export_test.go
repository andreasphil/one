package web

// Aliases and shims that let the external test package exercise the unexported
// mappers and the markdown renderer without exporting any of them to importers
// of this package.

type (
	NoteMeta         = noteMeta
	SearchResult     = searchResult
	MarkdownRenderer = markdownRenderer
)

var (
	NewNoteMeta         = newNoteMeta
	MapToNoteMeta       = mapToNoteMeta
	MapToTags           = mapToTags
	NewMarkdownRenderer = newMarkdownRenderer
	NewSearchResult     = newSearchResult
	MapToSearchResults  = mapToSearchResults
)

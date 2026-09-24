package cli

// Aliases and shims that let the external test package exercise the unexported
// notes provider without exporting it to importers of this package.

type FileNotesProvider = fileNotesProvider

var NewFileNotesProvider = newFileNotesProvider

package model

type Directory struct {
	Entries []DirectoryEntry
}

type DirectoryEntry struct {
	Name string
}

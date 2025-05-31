package model

type ContentType int

const (
	EmptyContentType     ContentType = iota
	DirectoryContentType ContentType = iota
	TextContentType      ContentType = iota
)

type Content interface {
	Type() ContentType
}

type EmptyContent struct{}

func (e *EmptyContent) Type() ContentType {
	return EmptyContentType
}

type DirectoryContent struct {
	Directory  Directory
	SelectedId int
}

func (d *DirectoryContent) Type() ContentType {
	return DirectoryContentType
}

type TextContent struct {
	Text string
}

func (t *TextContent) Type() ContentType {
	return TextContentType
}

type State struct {
	PrevViewContent Content
	CurrViewContent Content
	NextViewContent Content
	CurrentPath     string
}

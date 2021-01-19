package types

type ShareType int

const (
	InvalidType ShareType = iota
	TextType
	FileType
	// Directory
)

const (
	OK = 0
)

type Share struct {
	// for response
	Code int
	Msg  string

	Type ShareType `json:"type"`

	// text
	Content string `json:"content,omitempty"`

	// file
	FileName    string `json:"filename,omitempty"`
	FileContent []byte `json:"file_content,omitempty"`
}

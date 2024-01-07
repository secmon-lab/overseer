package model

type Alert struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Attrs       Attributes `json:"attrs"`
}

type AttrType string

const (
	IPAddr     AttrType = "ipaddr"
	DomainName AttrType = "domain"
	FileSha256 AttrType = "file.sha256"
	FileSha512 AttrType = "file.sha512"
	MarkDown   AttrType = "markdown"
)

type Attribute struct {
	Key   string   `json:"key" firestore:"key"`
	Value any      `json:"value" firestore:"value"`
	Type  AttrType `json:"type" firestore:"type"`
}

type Attributes []Attribute

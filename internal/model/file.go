package model

type SearchFile struct {
	Id   int
	Name string
}

type File struct {
	Id        int     `json:"id"`
	Url       string  `json:"url,omitempty"`
	PublicUrl string  `json:"public_url"`
	Name      string  `json:"name"`
	Size      int64   `json:"size"`
	MimeType  string  `json:"mime"`
	ViewName  *string `json:"view_name"`
}

func (f *File) GetViewName() string {
	if f.ViewName != nil {
		return *f.ViewName
	}

	return f.Name
}

package rssBodies

type Attr struct {
	Href string `json:"href"`
	Rel string `json:"rel"`
	Type string `json:"type"`
}

type Child struct {}

type LinksWithAttrs struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Attrs Attr `json:"attrs"`

	// seems that it is empty
	Children Child `json:"children"`
}

type Link struct {
	Link []LinksWithAttrs `json:"link"`
}

type AtomExtension struct {
	Atom Link `json:"atom"`
}

type Item struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`

	Published string `json:"published"`
	PublishedParsed string `json:"publishedParsed"`

	Guid string `json:"guid"`
	Categories []string  `json:"categories"`
}

type ITunesRSS struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`

	Updated string `json:"updated"`
	UpdatedParsed string `json:"updatedParsed"`
	Copyright string `json:"copyright"`
	Categories []string `json:"categories"`

	Extensions AtomExtension `json:"extensions"`

	Items []Item `json:"items"`

	FeedType string `json:"feedType"`
	FeedVersion string `json:"feedVersion"`


}

package pocket

type DomainMetadata struct {
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	GreyscaleLogo string `json:"greyscale_logo"`
}

type Item struct {
	ItemID              string         `json:"item_id"`
	NormalURL           string         `json:"normal_url"`
	ResolvedID          string         `json:"resolved_id"`
	ExtendedItemID      string         `json:"extended_item_id"`
	ResolvedURL         string         `json:"resolved_url"`
	DomainID            string         `json:"domain_id"`
	OriginDomainID      string         `json:"origin_domain_id"`
	ResponseCode        string         `json:"response_code"`
	MimeType            string         `json:"mime_type"`
	ContentLength       string         `json:"content_length"`
	Encoding            string         `json:"encoding"`
	DateResolved        string         `json:"date_resolved"`
	DatePublished       string         `json:"date_published"`
	Title               string         `json:"title"`
	Excerpt             string         `json:"excerpt"`
	WordCount           string         `json:"word_count"`
	InnerdomainRedirect string         `json:"innerdomain_redirect"`
	LoginRequired       string         `json:"login_required"`
	HasImage            string         `json:"has_image"`
	HasVideo            string         `json:"has_video"`
	IsIndex             string         `json:"is_index"`
	IsArticle           string         `json:"is_article"`
	UsedFallback        string         `json:"used_fallback"`
	Lang                string         `json:"lang"`
	TimeFirstParsed     string         `json:"time_first_parsed"`
	Authors             []interface{}  `json:"authors"`
	Images              []interface{}  `json:"images"`
	Videos              []interface{}  `json:"videos"`
	ResolvedNormalURL   string         `json:"resolved_normal_url"`
	DomainMetadata      DomainMetadata `json:"domain_metadata"`
	TimeToRead          int            `json:"time_to_read"`
	GivenURL            string         `json:"given_url"`
}

type AddResponse struct {
	Item   Item `json:"item"`
	Status int  `json:"status"`
}

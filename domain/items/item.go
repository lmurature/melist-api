package items

type Item struct {
	Id                string            `json:"id"`
	Title             string            `json:"title"`
	Descriptions      []ItemDescription `json:"descriptions"`
	CategoryId        string            `json:"category_id"`
	SellerId          int64             `json:"seller_id"`
	Price             float32           `json:"price"`
	Status            string            `json:"status"`
	InitialQuantity   int               `json:"initial_quantity"`
	AvailableQuantity int               `json:"available_quantity"`
	Condition         string            `json:"condition"`
	SoldQuantity      int               `json:"sold_quantity"`
	Pictures          []ItemPicture     `json:"pictures,omitempty"`
	Attributes        []ItemAttribute   `json:"attributes,omitempty"`
	SaleTerms         []ItemAttribute   `json:"sale_terms,omitempty"`
	Variations        []ItemVariation   `json:"variations,omitempty"`
	SubStatus         []string          `json:"sub_status"`
	Permalink         string            `json:"permalink,omitempty"`
}

type ItemDescription struct {
	Id        string `json:"id,omitempty"`
	Text      string `json:"text,omitempty"`
	PlainText string `json:"plain_text,omitempty"`
}

type ItemPicture struct {
	Id        string `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	SecureUrl string `json:"secure_url,omitempty"`
	Size      string `json:"size,omitempty"`
	MaxSize   string `json:"max_size,omitempty"`
}

type ItemAttribute struct {
	Id          string       `json:"id"`
	Name        string       `json:"name,omitempty"`
	ValueId     string       `json:"value_id,omitempty"`
	ValueName   string       `json:"value_name,omitempty"`
	ValueStruct *ValueStruct `json:"value_struct,omitempty"`
}

type ItemVariation struct {
	Id                    int64           `json:"id,omitempty"`
	AvailableQuantity     int             `json:"available_quantity"`
	Price                 float32         `json:"price"`
	Attributes            []ItemAttribute `json:"attributes,omitempty"`
	PictureIds            []string        `json:"picture_ids,omitempty"`
	AttributeCombinations []ItemAttribute `json:"attribute_combinations"`
}

type ValueStruct struct {
	Number float32 `json:"number"`
	Unit   string  `json:"unit"`
}

type ItemSearchResponse struct {
	SiteId         string `json:"site_id"`
	Query          string `json:"query"`
	Paging         Paging `json:"paging"`
	Result         []Item `json:"results"`
	Sort           Sort   `json:"sort"`
	AvailableSorts []Sort `json:"available_sorts"`
}

type Paging struct {
	Total  int64 `json:"total"`
	Offset int64 `json:"offset"`
	Limit  int64 `json:limit`
}

type Sort struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

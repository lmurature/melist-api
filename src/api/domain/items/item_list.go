package items

const (
	StatusChecked    = "checked"
	StatusNotChecked = "not_checked"
)

type ItemListDto struct {
	ItemId      string `json:"item_id"`
	ListId      int64  `json:"list_id"`
	Status      string `json:"status"`
	VariationId int64  `json:"variation_id,omitempty"`
	MeliItem    *Item  `json:"item,omitempty"`
}

type ItemListCollection []ItemListDto

func (items ItemListCollection) ContainsItem(itemId string) bool {
	for _, i := range items {
		if i.ItemId == itemId {
			return true
		}
	}
	return false
}

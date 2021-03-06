package items_provider

import (
	"github.com/lmurature/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "/sites/MLA/search?q=%s&offset=%d", uriSearchItems)
	assert.EqualValues(t, "/items/%s", uriGetItem)
}

func TestSearchItemsInvalidRestClientResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/sites/MLA/search?q=Computadora&offset=0",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: -1,
	})

	searchResult, err := SearchItemsByQuery("Computadora", 0)

	assert.Nil(t, searchResult)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response searching items", err.Message())
}

func TestSearchItemsInvalidApiErrorResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/sites/MLA/search?q=Computadora&offset=0&offset=0",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     `{---`,
	})

	searchResult, err := SearchItemsByQuery("Computadora", 0)

	assert.Nil(t, searchResult)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal apierror items search response", err.Message())
}

func TestSearchItemsErrorFromBytes(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/sites/MLA/search?q=Computadora&offset=0",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     `{"message": "internal server error trying to search items", "status": 500}`,
	})

	searchResult, err := SearchItemsByQuery("Computadora", 0)

	assert.Nil(t, searchResult)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "internal server error trying to search items", err.Message())
}

func TestSearchItemsErrorUnmarshal(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/sites/MLA/search?q=Computadora&offset=0",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{---`,
	})

	searchResult, err := SearchItemsByQuery("Computadora", 0)

	assert.Nil(t, searchResult)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal items search response", err.Message())
}

func TestSearchItemsNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/sites/MLA/search?q=Computadora&offset=0",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"site_id":"MLA","query":"Computadora","paging":{"total":20027,"offset":0,"limit":50},"results":[{"id":"MLA907751590","title":"Estabilizador De Tensi??n Lyonn Tca Series 1200nv 1200va Con Entrada Y Salida De 220v Ca  Negro","descriptions":null,"category_id":"MLA1719","seller_id":0,"price":2564,"status":"","initial_quantity":0,"available_quantity":11,"condition":"new","sold_quantity":314,"attributes":[{"id":"BRAND","name":"Marca","value_id":"15747","value_name":"Lyonn"},{"id":"ITEM_CONDITION","name":"Condici??n del ??tem","value_id":"2230284","value_name":"Nuevo"},{"id":"LINE","name":"L??nea","value_id":"338326","value_name":"TCA Series"},{"id":"MODEL","name":"Modelo","value_id":"9729806","value_name":"1200NV"},{"id":"PEAK_POWER","name":"Potencia pico","value_id":"260601","value_name":"1200VA","value_struct":{"number":1200,"unit":"VA"}},{"id":"RATED_POWER","name":"Potencia nominal","value_id":"8900723","value_name":"1200 VA","value_struct":{"number":1200,"unit":"VA"}},{"id":"WEIGHT","name":"Peso","value_id":"7726408","value_name":"1.54 kg","value_struct":{"number":1.54,"unit":"kg"}}],"sub_status":null,"permalink":"https://www.mercadolibre.com.ar/estabilizador-de-tension-lyonn-tca-series-1200nv-1200va-con-entrada-y-salida-de-220v-ca-negro/p/MLA6208662"},{"id":"MLA873398163","title":"Memoria Ram Fury Ddr4 Gamer 8gb 1x8gb Hyperx Hx426c16fb3/8","descriptions":null,"category_id":"MLA1694","seller_id":0,"price":6319,"status":"","initial_quantity":0,"available_quantity":2672,"condition":"new","sold_quantity":4341,"attributes":[{"id":"BRAND","name":"Marca","value_id":"448156","value_name":"HyperX"},{"id":"ITEM_CONDITION","name":"Condici??n del ??tem","value_id":"2230284","value_name":"Nuevo"},{"id":"LINE","name":"L??nea","value_id":"10087029","value_name":"Fury DDR4"},{"id":"MODEL","name":"Modelo","value_id":"7790422","value_name":"HX426C16FB3/8"},{"id":"PACKAGE_LENGTH","name":"Largo del paquete","value_name":"13.6 cm","value_struct":{"number":13.6,"unit":"cm"}},{"id":"PACKAGE_WEIGHT","name":"Peso del paquete","value_name":"60 g","value_struct":{"number":60,"unit":"g"}}],"sub_status":null,"permalink":"https://www.mercadolibre.com.ar/memoria-ram-fury-ddr4-gamer-8gb-1x8gb-hyperx-hx426c16fb38/p/MLA15178125"},{"id":"MLA879276614","title":"Computadora Cpu Intel Amd Doble Nucleo 8 Gb 500 Gb","descriptions":null,"category_id":"MLA1649","seller_id":0,"price":26590,"status":"","initial_quantity":0,"available_quantity":1,"condition":"new","sold_quantity":150,"attributes":[{"id":"BRAND","name":"Marca","value_id":"18034","value_name":"AMD"},{"id":"ITEM_CONDITION","name":"Condici??n del ??tem","value_id":"2230284","value_name":"Nuevo"},{"id":"MODEL","name":"Modelo","value_name":"AMD E6010"},{"id":"PACKAGE_LENGTH","name":"Largo del paquete","value_name":"45.8 cm","value_struct":{"number":45.8,"unit":"cm"}},{"id":"PACKAGE_WEIGHT","name":"Peso del paquete","value_name":"5880 g","value_struct":{"number":5880,"unit":"g"}}],"sub_status":null,"permalink":"https://articulo.mercadolibre.com.ar/MLA-879276614-computadora-cpu-intel-amd-doble-nucleo-8-gb-500-gb-_JM"}],"sort":{"id":"relevance","name":"M??s relevantes"},"available_sorts":[{"id":"price_asc","name":"Menor precio"},{"id":"price_desc","name":"Mayor precio"}]}`,
	})

	searchResult, err := SearchItemsByQuery("Computadora", 0)

	assert.Nil(t, err)
	assert.NotNil(t, searchResult)
	assert.EqualValues(t, "MLA", searchResult.SiteId)
	assert.EqualValues(t, "Computadora", searchResult.Query)
	assert.EqualValues(t, 20027, searchResult.Paging.Total)
	assert.EqualValues(t, 0, searchResult.Paging.Offset)
	assert.EqualValues(t, 50, searchResult.Paging.Limit)
	assert.EqualValues(t, 3, len(searchResult.Result))
	assert.EqualValues(t, "MLA907751590", searchResult.Result[0].Id)
}

func TestGetItemInvalidRestClientResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/items/MLA1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: -1,
	})

	item, err := GetItemById("MLA1")
	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response getting item MLA1", err.Message())
}

func TestGetItemInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/items/MLA1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "Item with id MLA1 not found.","error": "not_found", "status": "404", "cause": []}`,
	})

	item, err := GetItemById("MLA1")
	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal apierror item response", err.Message())
}

func TestGetItemNotFound(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/items/MLA1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "Item with id MLA1 not found.","error": "not_found", "status": 404, "cause": []}`,
	})

	item, err := GetItemById("MLA1")
	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "Item with id MLA1 not found.", err.Message())
}

func TestGetItemInvalidJson(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/items/MLA1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{----`,
	})

	item, err := GetItemById("MLA1")
	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error trying to unmarshal response into item MLA1 structure", err.Message())
}

func TestGetItemNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/items/MLA1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":"MLA1","site_id":"MLA","title":"Test item - DO NOT BUY","descriptions":[{"plain_text": "this is the description"}],"listing_type_id":"gold_pro","category_id":"CBT412445","seller_id":460986913,"price":500,"base_price":500,"initial_quantity":10,"available_quantity":9,"sold_quantity":1, "status": "active"}`,
	})

	item, err := GetItemById("MLA1")
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.EqualValues(t, "MLA1", item.Id)
	assert.EqualValues(t, "Test item - DO NOT BUY", item.Title)
	assert.EqualValues(t, "CBT412445", item.CategoryId)
	assert.EqualValues(t, 460986913, item.SellerId)
	assert.EqualValues(t, 500, item.Price)
	assert.EqualValues(t, "active", item.Status)
	assert.EqualValues(t, 10, item.InitialQuantity)
	assert.EqualValues(t, 9, item.AvailableQuantity)
	assert.EqualValues(t, 1, item.SoldQuantity)
}

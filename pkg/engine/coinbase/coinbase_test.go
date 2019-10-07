package coinbase

import (
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestToOrder(t *testing.T) {
	r := require.New(t)
	rawOrderLV2 := []interface{}{"1", "2", 3}
	odLV2, err := ToOrder(rawOrderLV2)
	r.NoError(err)
	r.Equal(decimal.RequireFromString("1"), odLV2.Price)
	r.Equal(decimal.RequireFromString("2"), odLV2.Size)
	r.Equal("", odLV2.OrderID)
	r.Equal(int64(3), odLV2.NumOrders)

	rawOrderLV3 := []interface{}{"295.97", "5.72036512", "da863862-25f4-4868-ac41-005d11ab0a5f"}
	odLV3, err := ToOrder(rawOrderLV3)
	r.NoError(err)
	r.Equal(decimal.RequireFromString("295.97"), odLV3.Price)
	r.Equal(decimal.RequireFromString("5.72036512"), odLV3.Size)
	r.Equal("da863862-25f4-4868-ac41-005d11ab0a5f", odLV3.OrderID)
	r.Equal(int64(0), odLV3.NumOrders)
}

func TestToOrderBook(t *testing.T) {
	r := require.New(t)
	raw := []byte(`{
		"sequence":7371656227,
		"bids": [["170.95","0.34084807",1]],
		"asks": [["170.97","7.64562173",3]]
	}`)

	rawLV2 := []byte(`{"sequence":7371989985,"bids":[["170.22","0.12421147",1],["170.2","27.49532233",1],["170.19","27.83392529",1],["170.16","28.15088947",2],["170.12","7.4703393",1],["170.08","41.75393453",2],["170.06","110",1],["170.05","50.7",1],["170.03","2.23",1],["170.01","30.81022191",2],["170","4",1],["169.97","44.08",1],["169.95","45.36337329",2],["169.94","5",1],["169.93","262.9",2],["169.85","40.9",1],["169.84","24.34949871",1],["169.82","25",1],["169.75","2.91",1],["169.74","31",1],["169.69","16.13035732",3],["169.64","174",1],["169.63","34.78",1],["169.61","1.17515717",1],["169.6","44.30784543",3],["169.52","63.77576864",2],["169.5","17.028",2],["169.49","39.4",1],["169.45","63.06",1],["169.42","48.81",2],["169.41","1",1],["169.4","21.51756472",1],["169.33","19.39",1],["169.31","38.10224349",2],["169.18","5",1],["169.14","63.62",1],["169.09","3.55939201",1],["169.05","41.37826365",1],["169","41.10499079",1],["168.99","4.95",1],["168.8","0.52777",1],["168.78","100",1],["168.73","632.79",1],["168.63","34.198",1],["168.41","1",1],["168.3","0.5",1],["168.08","0.02",1],["168.06","150.97884684",1],["168.05","40",1],["168","152.61681752",5]],"asks":[["170.23","13.32546631",2],["170.24","0.28124951",1],["170.26","46.833",1],["170.27","10",1],["170.32","110",1],["170.33","174",1],["170.34","2",1],["170.38","3",1],["170.43","20.27014126",1],["170.44","224",2],["170.45","5",1],["170.48","49.43924335",2],["170.5","18.36",1],["170.59","27.49515791",2],["170.6","34.91999999",2],["170.61","360",1],["170.62","329.01474036",1],["170.63","8",1],["170.67","57.03841",2],["170.71","30.9768506",2],["170.75","25.67383096",2],["170.76","0.017",1],["170.78","23.42331",1],["170.79","110.65",2],["170.81","38.2",1],["170.82","0.17685453",1],["170.83","0.67688635",2],["170.84","9.91305888",1],["170.85","21.49424265",1],["170.88","0.5",1],["170.89","0.5",1],["170.91","20.1",2],["170.93","0.017",1],["170.96","23.93170437",1],["170.97","31.1",1],["171","93.0595235",3],["171.04","40.092",2],["171.08","22.7074476",1],["171.12","1",1],["171.15","16",1],["171.17","20.81905893",1],["171.31","63.61",1],["171.34","0.5",1],["171.55","21.807",1],["171.65","5.89",1],["171.74","0.14",1],["171.77","0.5",1],["171.79","0.5",1],["171.91","603.49475",3],["171.92","632.79",1]]}`)

	rawLV3 := []byte(`{"sequence":7371931270,"bids":[["170.26","1","c443f5ce-354d-43e6-a279-a0cea25e3f32"],["170.26","10","451a73da-9148-491c-83b5-eb37ba9ed510"],["170.24","4","07c17105-d5a2-4d2d-9867-825480156cbe"],["170.17","44","f96dda3c-1ba1-4bee-b42b-4298d6448ce9"],["170.15","10","cb50bf05-eccf-4da0-bfe2-b269a49c8d90"],["170.12","0.02","18c3e808-9a7d-4d0c-9eea-7b1b6c6295ed"],["170.12","110","b82de571-64c7-48f7-8a0f-6d10825a1d77"],["170.11","43.97","43bf118b-b862-47d6-af9a-577da129525f"],["170.11","10","8b2d95c8-f915-4fe7-a29b-865a34f3bc47"],["170.08","7","ef75fa65-0e24-4262-881c-6f22949b89bc"],["170.05","20.16205373","f99d5356-a1cc-464c-b5f6-021b85f72835"],["170.05","46.722","c9971e40-1c96-4af4-982c-6212f0881a62"],["170.04","7.4703393","124e5667-55e1-4b69-9f47-49b0a1fba250"],["170","20","51f692e5-e89c-4214-b5f8-284d3b0c685f"],["169.99","38.872","db65376e-c373-45af-b0fb-cd2acd27e144"],["169.99","13.05","4afd4081-9cef-4a70-9ca1-e835d657d745"]],"asks":[["170.29","21.128","dd98f6ea-6e2b-496d-971b-81df71cbb66e"],["170.32","4","7752dac7-0345-4b06-a01a-8861844f6e6d"],["170.41","20.16148875","b5f8ee42-3558-4d71-b5dc-4ba5ac5e4109"],["170.41","174","d13d084a-f380-4df0-9cfc-6b7160147872"],["170.51","25","f29abe4b-a966-484a-b482-c9974f522a60"],["170.51","0.07222","1b68f3c6-3270-4d71-be71-1fe8870c4313"],["170.53","360","0df8f54b-f2fb-4076-ae05-9ce0f1a695c4"],["170.54","31.1","f4779bed-24ab-402b-b3f2-11962bd62499"],["170.54","5","990736de-0393-442c-b941-66f2ffc29855"]]}`)

	lv1, err := ToOrderBook(raw, time.Now())
	r.NoError(err)
	r.Len(lv1.Bids, 1)
	r.Len(lv1.Asks, 1)
	r.Equal(lv1.Sequence, "7371656227")
	for _, order := range lv1.Bids {
		r.Empty(order.OrderID)
		r.NotEmpty(order.NumOrders)
	}
	for _, order := range lv1.Asks {
		r.Empty(order.OrderID)
		r.NotEmpty(order.NumOrders)
	}

	lv2, err := ToOrderBook(rawLV2, time.Now())
	r.NoError(err)
	r.Len(lv2.Bids, 50)
	r.Len(lv2.Asks, 50)
	r.Equal(lv2.Sequence, "7371989985")
	for _, order := range lv1.Bids {
		r.Empty(order.OrderID)
		r.NotEmpty(order.NumOrders)
	}
	for _, order := range lv1.Asks {
		r.Empty(order.OrderID)
		r.NotEmpty(order.NumOrders)
	}

	lv3, err := ToOrderBook(rawLV3, time.Now())
	r.Len(lv3.Bids, 16)
	r.Len(lv3.Asks, 9)
	for _, order := range lv3.Bids {
		r.NotEmpty(order.OrderID)
		r.Empty(order.NumOrders)
	}
	for _, order := range lv3.Asks {
		r.NotEmpty(order.OrderID)
		r.Empty(order.NumOrders)
	}
}

// this test actually make request to network
// skip when testing in development version
func SkipTestFetchOrderBook(t *testing.T) {
	coinbaseAPIURL := "https://api.pro.coinbase.com"
	tests := []struct {
		endpoint  string
		level     int64
		pair      string
		wantError bool
	}{
		{
			endpoint:  coinbaseAPIURL,
			level:     1,
			pair:      "ETH-USD",
			wantError: false,
		},
		{
			endpoint:  coinbaseAPIURL,
			level:     2,
			pair:      "BTC-USD",
			wantError: false,
		},
		{
			endpoint:  coinbaseAPIURL,
			level:     3,
			pair:      "BTC-USD",
			wantError: true,
		},
	}

	r := require.New(t)

	for _, tt := range tests {
		resp, _, err := FetchOrderBook(tt.endpoint, tt.level, tt.pair)
		if tt.wantError {
			r.Error(err)
		} else {
			r.NoError(err)
			r.NotNil(resp)
			JSONErr := validation.Validate(resp, is.JSON)
			r.NoError(JSONErr)
		}
	}
}

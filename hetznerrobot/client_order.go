package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#server-ordering

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type HetznerRobotProductPrice struct {
	Location     string
	MonthlyNet   string
	MonthlyGross string
	SetupNet     string
	SetupGross   string
}

type HetznerRobotServerProduct struct {
	ID          string
	Name        string
	Description []string
	Traffic     string
	Locations   []string
	Prices      []HetznerRobotProductPrice
	// orderable options (populated by the product-detail endpoint)
	Dist []string
	Lang []string
	Arch []string
}

type HetznerRobotServerTransaction struct {
	ID           string
	Status       string
	ServerNumber int
	ServerIP     string
}

type HetznerRobotServerOrderRequest struct {
	ProductID      string
	Location       string
	Dist           string
	Lang           string
	AuthorizedKeys []string
	Password       string
	Addons         []string
	Test           bool
}

func gjsonStrings(r gjson.Result) []string {
	out := []string{}
	r.ForEach(func(_, v gjson.Result) bool {
		out = append(out, v.String())
		return true
	})
	return out
}

func parseServerProduct(p gjson.Result) HetznerRobotServerProduct {
	prices := []HetznerRobotProductPrice{}
	p.Get("prices").ForEach(func(_, pr gjson.Result) bool {
		prices = append(prices, HetznerRobotProductPrice{
			Location:     pr.Get("location").String(),
			MonthlyNet:   pr.Get("price.net").String(),
			MonthlyGross: pr.Get("price.gross").String(),
			SetupNet:     pr.Get("price_setup.net").String(),
			SetupGross:   pr.Get("price_setup.gross").String(),
		})
		return true
	})
	return HetznerRobotServerProduct{
		ID:          p.Get("id").String(),
		Name:        p.Get("name").String(),
		Description: gjsonStrings(p.Get("description")),
		Traffic:     p.Get("traffic").String(),
		Locations:   gjsonStrings(p.Get("location")),
		Prices:      prices,
		Dist:        gjsonStrings(p.Get("dist")),
		Lang:        gjsonStrings(p.Get("lang")),
		Arch:        gjsonStrings(p.Get("arch")),
	}
}

func parseServerTransaction(jsonStr string) *HetznerRobotServerTransaction {
	t := gjson.Get(jsonStr, "transaction")
	if !t.Exists() {
		t = gjson.Parse(jsonStr)
	}
	return &HetznerRobotServerTransaction{
		ID:           t.Get("id").String(),
		Status:       t.Get("status").String(),
		ServerNumber: int(t.Get("server_number").Int()),
		ServerIP:     t.Get("server_ip").String(),
	}
}

func (c *HetznerRobotClient) getServerProducts(ctx context.Context) ([]HetznerRobotServerProduct, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server/product", c.url), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	products := []HetznerRobotServerProduct{}
	gjson.ParseBytes(res).ForEach(func(_, item gjson.Result) bool {
		p := item.Get("product")
		if !p.Exists() {
			p = item
		}
		products = append(products, parseServerProduct(p))
		return true
	})
	return products, nil
}

// getServerProduct returns a single product's details including its orderable
// dist/lang/arch/location options (GET /order/server/product/{id}).
func (c *HetznerRobotClient) getServerProduct(ctx context.Context, id string) (*HetznerRobotServerProduct, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server/product/%s", c.url, url.PathEscape(id)), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	p := gjson.GetBytes(res, "product")
	if !p.Exists() {
		p = gjson.ParseBytes(res)
	}
	product := parseServerProduct(p)
	return &product, nil
}

// createServerOrder places a server order transaction. When req.Test is true the
// order is validated only (no server is provisioned and no charge is made).
func (c *HetznerRobotClient) createServerOrder(ctx context.Context, req HetznerRobotServerOrderRequest) (*HetznerRobotServerTransaction, error) {
	data := url.Values{}
	data.Set("product_id", req.ProductID)
	if req.Location != "" {
		data.Set("location", req.Location)
	}
	if req.Dist != "" {
		data.Set("dist", req.Dist)
	}
	if req.Lang != "" {
		data.Set("lang", req.Lang)
	}
	for _, key := range req.AuthorizedKeys {
		data.Add("authorized_key[]", key)
	}
	if req.Password != "" {
		data.Set("password", req.Password)
	}
	for _, addon := range req.Addons {
		data.Add("addon[]", addon)
	}
	if req.Test {
		data.Set("test", "true")
	}

	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/order/server/transaction", c.url), data, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}

func (c *HetznerRobotClient) getServerOrder(ctx context.Context, id string) (*HetznerRobotServerTransaction, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server/transaction/%s", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}

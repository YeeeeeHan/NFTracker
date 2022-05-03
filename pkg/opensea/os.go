package opensea

import (
	"NFTracker/pkg/custError"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"log"
	"net/http"
)

type OSResponse struct {
	Collection struct {
		Editors       []string `json:"editors"`
		PaymentTokens []struct {
			ID       int     `json:"id"`
			Symbol   string  `json:"symbol"`
			Address  string  `json:"address"`
			ImageURL string  `json:"image_url"`
			Name     string  `json:"name"`
			Decimals int     `json:"decimals"`
			EthPrice float64 `json:"eth_price"`
			UsdPrice float64 `json:"usd_price"`
		} `json:"payment_tokens"`
		PrimaryAssetContracts []struct {
			Address                     string      `json:"address"`
			AssetContractType           string      `json:"asset_contract_type"`
			CreatedDate                 string      `json:"created_date"`
			Name                        string      `json:"name"`
			NftVersion                  string      `json:"nft_version"`
			OpenseaVersion              interface{} `json:"opensea_version"`
			Owner                       int         `json:"owner"`
			SchemaName                  string      `json:"schema_name"`
			Symbol                      string      `json:"symbol"`
			TotalSupply                 string      `json:"total_supply"`
			Description                 string      `json:"description"`
			ExternalLink                string      `json:"external_link"`
			ImageURL                    string      `json:"image_url"`
			DefaultToFiat               bool        `json:"default_to_fiat"`
			DevBuyerFeeBasisPoints      int         `json:"dev_buyer_fee_basis_points"`
			DevSellerFeeBasisPoints     int         `json:"dev_seller_fee_basis_points"`
			OnlyProxiedTransfers        bool        `json:"only_proxied_transfers"`
			OpenseaBuyerFeeBasisPoints  int         `json:"opensea_buyer_fee_basis_points"`
			OpenseaSellerFeeBasisPoints int         `json:"opensea_seller_fee_basis_points"`
			BuyerFeeBasisPoints         int         `json:"buyer_fee_basis_points"`
			SellerFeeBasisPoints        int         `json:"seller_fee_basis_points"`
			PayoutAddress               string      `json:"payout_address"`
		} `json:"primary_asset_contracts"`
		Stats struct {
			OneDayVolume          float64 `json:"one_day_volume"`
			OneDayChange          float64 `json:"one_day_change"`
			OneDaySales           float64 `json:"one_day_sales"`
			OneDayAveragePrice    float64 `json:"one_day_average_price"`
			SevenDayVolume        float64 `json:"seven_day_volume"`
			SevenDayChange        float64 `json:"seven_day_change"`
			SevenDaySales         float64 `json:"seven_day_sales"`
			SevenDayAveragePrice  float64 `json:"seven_day_average_price"`
			ThirtyDayVolume       float64 `json:"thirty_day_volume"`
			ThirtyDayChange       float64 `json:"thirty_day_change"`
			ThirtyDaySales        float64 `json:"thirty_day_sales"`
			ThirtyDayAveragePrice float64 `json:"thirty_day_average_price"`
			TotalVolume           float64 `json:"total_volume"`
			TotalSales            float64 `json:"total_sales"`
			TotalSupply           float64 `json:"total_supply"`
			Count                 float64 `json:"count"`
			NumOwners             int     `json:"num_owners"`
			AveragePrice          float64 `json:"average_price"`
			NumReports            int     `json:"num_reports"`
			MarketCap             float64 `json:"market_cap"`
			FloorPrice            float64 `json:"floor_price"`
		} `json:"stats"`
		BannerImageURL          string      `json:"banner_image_url"`
		ChatURL                 interface{} `json:"chat_url"`
		CreatedDate             string      `json:"created_date"`
		DefaultToFiat           bool        `json:"default_to_fiat"`
		Description             string      `json:"description"`
		DevBuyerFeeBasisPoints  string      `json:"dev_buyer_fee_basis_points"`
		DevSellerFeeBasisPoints string      `json:"dev_seller_fee_basis_points"`
		DiscordURL              string      `json:"discord_url"`
		DisplayData             struct {
			CardDisplayStyle string `json:"card_display_style"`
		} `json:"display_data"`
		ExternalURL                 string      `json:"external_url"`
		Featured                    bool        `json:"featured"`
		FeaturedImageURL            interface{} `json:"featured_image_url"`
		Hidden                      bool        `json:"hidden"`
		SafelistRequestStatus       string      `json:"safelist_request_status"`
		ImageURL                    string      `json:"image_url"`
		IsSubjectToWhitelist        bool        `json:"is_subject_to_whitelist"`
		LargeImageURL               interface{} `json:"large_image_url"`
		MediumUsername              interface{} `json:"medium_username"`
		Name                        string      `json:"name"`
		OnlyProxiedTransfers        bool        `json:"only_proxied_transfers"`
		OpenseaBuyerFeeBasisPoints  string      `json:"opensea_buyer_fee_basis_points"`
		OpenseaSellerFeeBasisPoints string      `json:"opensea_seller_fee_basis_points"`
		PayoutAddress               string      `json:"payout_address"`
		RequireEmail                bool        `json:"require_email"`
		ShortDescription            interface{} `json:"short_description"`
		Slug                        string      `json:"slug"`
		TelegramURL                 interface{} `json:"telegram_url"`
		TwitterUsername             interface{} `json:"twitter_username"`
		InstagramUsername           string      `json:"instagram_username"`
		WikiURL                     interface{} `json:"wiki_url"`
	} `json:"collection"`
}

func (osr OSResponse) GetSupplyString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(osr.Collection.Stats.TotalSupply))
}

func (osr OSResponse) GetFloorPriceString() string {
	return fmt.Sprintf("%.2f", osr.Collection.Stats.FloorPrice)
}

func (osr OSResponse) GetNumOwnersString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", osr.Collection.Stats.NumOwners)
}

func (osr OSResponse) GetOneDayChangeString() string {
	return fmt.Sprintf("%.3f", osr.Collection.Stats.OneDayChange)
}

func (osr OSResponse) GetOneDayVolumeString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(osr.Collection.Stats.OneDayVolume))
}

func (osr OSResponse) GetSevenDayChangeString() string {
	return fmt.Sprintf("%.3f", osr.Collection.Stats.SevenDayChange)
}

func (osr OSResponse) GetSevenDayVolumeString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(osr.Collection.Stats.SevenDayVolume))
}

func (osr OSResponse) GetTotalVolumeString() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(osr.Collection.Stats.TotalVolume))
}

func QueryAPI(slug string) (*OSResponse, error) {

	url := "https://api.opensea.io/collection/" + slug

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Printf("[http.DefaultClient.Do] Non-OK HTTP status: %d", res.StatusCode)
		return nil, custError.InvalidSlugErr
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	osResponse := OSResponse{}
	if err := json.Unmarshal(body, &osResponse); err != nil {
		log.Printf("[json.Unmarshal]: %v", err)
		return nil, err
	}

	return &osResponse, nil
}

func CreateOpenseaUrlFromSlug(slug string) string {
	return "https://opensea.io/collection/" + slug
}

func CreateTwitterUrlFromSlug(slug string) string {
	return "https://twitter.com/" + slug
}

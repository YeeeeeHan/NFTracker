package os

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Scrape(slug string) {
	url := "https://api.opensea.io/collection/" + slug

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	return
}

package cmd

import (
	"NFTracker/pkg/opensea"
	"fmt"
)

var WelcomeMessage = "Welcome to the 🌊 *OpenSea NFTracker Bot* 🌊\n\n" +
	"👉 Use `/check <slug>` to receive the latest stats of a collection _(slug is the collection name in the url, e.g. doodles-official in https://opensea.io/collection/doodles-official)_\n" +
	"👉 More features coming soon..." +
	"\n\nBot creator: [Yee Han](https://github.com/YeeeeeHan)"

func PriceCheckMessage(slug string, link string, osr *opensea.OSResponse) string {
	return fmt.Sprintf("[%s](%s)\n"+
		"✨ Floor price:                 %sΞ\n\n"+
		"📉 1-Day FP change:     %sΞ\n"+
		"📦 1-Day volume:          %sΞ\n"+
		"📈 7-Day FP change:     %sΞ\n"+
		"📦 7-Day volume:          %sΞ\n"+
		"💎 Volume traded:        %sΞ\n"+
		"💯 Supply:                       %s\n"+
		"✋🏼 Owners:                      %s",
		slug,
		link,
		osr.GetFloorPriceString(),
		osr.GetOneDayChangeString(),
		osr.GetOneDayVolumeString(),
		osr.GetSevenDayChangeString(),
		osr.GetSevenDayVolumeString(),
		osr.GetTotalVolumeString(),
		osr.GetSupplyString(),
		osr.GetNumOwnersString())
}

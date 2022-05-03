package message

import (
	"NFTracker/pkg/opensea"
	"fmt"
)

var WelcomeMessage = "Welcome to the   *OpenSea NFTracker Bot* \n\n" +
	"👉 Use `/fp <slug>` to receive the latest stats of a collection _(slug is the collection name in the url, e.g. doodles-official in https://opensea.io/collection/doodles-official)_\n\n" +
	"👉 More features coming soon..."

func PriceCheckMessage(collectionName string, link string, osr *opensea.OSResponse) string {
	return fmt.Sprintf("[%s](%s)\n"+
		"✨ Floor price:               %sΞ\n\n"+
		"📦 1-Day volume:          %sΞ\n"+
		"📦 7-Day volume:          %sΞ\n"+
		"💎 Volume traded:        %sΞ\n"+
		"💯 Supply:                       %s\n"+
		"✋🏼 Owners:                      %s\n"+
		"🌊 [Visit Opensea](%s)\n"+
		"%s",
		collectionName,
		link,
		osr.GetFloorPriceString(),
		osr.GetOneDayVolumeString(),
		osr.GetSevenDayVolumeString(),
		osr.GetTotalVolumeString(),
		osr.GetSupplyString(),
		osr.GetNumOwnersString(),
		link,
		renderTwitterField(osr))
}

func renderTwitterField(osr *opensea.OSResponse) string {
	if osr.Collection.TwitterUsername == nil {
		return ""
	}

	return "🐦 [Visit Twitter](" + opensea.CreateTwitterUrlFromSlug(fmt.Sprintf("%v", osr.Collection.TwitterUsername)) + ")"
}

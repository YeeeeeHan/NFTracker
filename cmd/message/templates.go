package message

import (
	"NFTracker/pkg/opensea"
	"fmt"
)

var WelcomeMessage = "Welcome to the   *OpenSea NFTracker Bot* \n\n" +
	"š Use `/fp <slug>` to receive the latest stats of a collection _(slug is the collection name in the url, e.g. doodles-official in https://opensea.io/collection/doodles-official)_\n\n" +
	"š More features coming soon..."

func PriceCheckMessage(collectionName string, link string, osr *opensea.OSResponse) string {
	return fmt.Sprintf("[%s](%s)\n"+
		"āØ Floor price:               %sĪ\n\n"+
		"š¦ 1-Day volume:          %sĪ\n"+
		"š¦ 7-Day volume:          %sĪ\n"+
		"š Volume traded:        %sĪ\n"+
		"šÆ Supply:                       %s\n"+
		"āš¼ Owners:                      %s\n"+
		"\n"+
		"š [Visit Opensea](%s)\n"+
		"%s\n"+
		"%s\n"+
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
		renderDiscordField(osr),
		renderLooksRareField(osr),
		renderTwitterField(osr))
}

func renderTwitterField(osr *opensea.OSResponse) string {
	if osr.Collection.TwitterUsername == nil {
		return ""
	}
	return "š¦ [Visit Twitter](" + opensea.CreateTwitterUrlFromUsername(osr.GetTwitterUsername()) + ")"
}

func renderDiscordField(osr *opensea.OSResponse) string {
	if osr.GetDiscordURL() == "" {
		return ""
	}
	return "š¾ [Visit Discord](" + opensea.CreateTwitterUrlFromUsername(osr.GetDiscordURL()) + ")"
}

func renderLooksRareField(osr *opensea.OSResponse) string {
	ca := osr.GetContractAddress()
	if ca == "" {
		return ""
	}

	return "āļø [Visit LooksRare](" + opensea.CreateLooksrareUrlFromAddress(ca) + ")"
}

package services

import (
	"fmt"
	"lxm-ency/utils"
	"net/http"
)

func encyCardHandle(w http.ResponseWriter, r *http.Request) {
	// get card
	cid := r.URL.Query().Get("cid")
	card, err := utils.GetCard(cid)
	if err != nil {
		returnErr(err, w)
		return
	}

	// return card
	writeCard(card, "dummy", w)
}

// output results
func writeCard(card *utils.Card, dummy string, w http.ResponseWriter) {
	resp := &Response{}
	cardWrapped := map[string]utf8String{
		"title":           utf8String(card.Title),
		"description":     utf8String(card.Description),
		"cover_url":       utf8String(card.CoverURL),
		"destination_url": utf8String(card.DestURL),
	}
	resp.Messages = append(resp.Messages, ResponseMsg{cardWrapped, "share_link"})
	resp.Slots = append(resp.Slots, Slots{"ency_dummy", dummy})
	fmt.Printf("[Response]%s \n", resp.ToBytes())
	w.Write(resp.ToBytes())
}

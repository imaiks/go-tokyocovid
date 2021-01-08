package main

import (
	"fmt"
	"sort"
)

var tokyoID = map[string]int{
	"千代田":  13101,
	"中央":   13102,
	"港":    13103,
	"新宿":   13104,
	"文京":   13105,
	"台東":   13106,
	"墨田":   13107,
	"江東":   13108,
	"品川":   13109,
	"目黒":   13110,
	"大田":   13111,
	"世田谷":  13112,
	"渋谷":   13113,
	"中野":   13114,
	"杉並":   13115,
	"豊島":   13116,
	"北":    13117,
	"荒川":   13118,
	"板橋":   13119,
	"練馬":   13120,
	"足立":   13121,
	"葛飾":   13122,
	"江戸川":  13123,
	"八王子":  13201,
	"立川":   13202,
	"武蔵野":  13203,
	"三鷹":   13204,
	"青梅":   13205,
	"府中":   13206,
	"昭島":   13207,
	"調布":   13208,
	"町田":   13209,
	"小金井":  13210,
	"小平":   13211,
	"日野":   13212,
	"東村山":  13213,
	"国分寺":  13214,
	"国立":   13215,
	"福生":   13218,
	"狛江":   13219,
	"東大和":  13220,
	"清瀬":   13221,
	"東久留米": 13222,
	"武蔵村山": 13223,
	"多摩":   13224,
	"稲城":   13225,
	"羽村":   13227,
	"あきる野": 13228,
	"西東京":  13229,
	"瑞穂":   13303,
	"日の出":  13305,
	"檜原":   13307,
	"奥多摩":  13308,
	"大島":   13361,
	"利島":   13362,
	"新島":   13363,
	"神津島":  13364,
	"三宅":   13381,
	"御蔵島":  13382,
	"八丈":   13401,
	"青ヶ島":  13402,
	"小笠原":  13421,
}

func TokyoJISCodes() (codes []string) {
	ids := []int{}
	for _, id := range tokyoID {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	for _, id := range ids {
		codes = append(codes, fmt.Sprintf("%d", id))
	}
	return
}
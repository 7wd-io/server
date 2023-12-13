package main

import (
	"encoding/json"
	"fmt"
	"github.com/7wd-io/engine"
	"strings"
)

type Item struct {
	Move interface{} `json:"move"`
}

func main() {
	raw := `[{"meta": {"actor": ""}, "move": {"id": 1, "p1": "user2", "p2": "user1", "cards": {"1": [112, 101, 110, 120, 114, 104, 121, 103, 100, 111, 118, 108, 122, 102, 113, 115, 116, 117, 119, 106], "2": [213, 218, 206, 200, 211, 217, 201, 205, 221, 203, 209, 212, 222, 214, 220, 202, 208, 210, 215, 207], "3": [318, 313, 307, 302, 309, 403, 303, 400, 301, 306, 304, 300, 315, 311, 308, 314, 316, 401, 312, 310]}, "tokens": [4, 1, 3, 10, 6], "wonders": [10, 5, 2, 14, 12, 8, 6, 7], "randomTokens": [7, 2, 5]}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 10}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 5}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 14}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 2}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 8}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 6}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 12}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 7}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 113}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 106}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 117}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 119}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 122}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 102}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 100}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 116}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 115}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 111}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 118}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 108}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 104}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 121}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 110}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 103}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 114}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 120}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 112}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 101}}, {"meta": {"actor": "user2"}, "move": {"id": 7, "player": "user2"}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 215}}, {"meta": {"actor": "user2"}, "move": {"id": 3, "token": 3}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 202}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 212}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 207}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 208}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 222}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 205}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 201}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 210}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 220}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 214}}, {"meta": {"actor": "user2"}, "move": {"id": 3, "token": 10}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 203, "wonder": 5}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 213}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 218}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 209, "wonder": 2}}, {"meta": {"actor": "user2"}, "move": {"id": 8, "card": 106}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 211}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 217}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 221, "wonder": 8}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 206}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 200}}, {"meta": {"actor": "user1"}, "move": {"id": 7, "player": "user1"}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 312}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 314}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 300, "wonder": 14}}, {"meta": {"actor": "user1"}, "move": {"id": 12, "give": 319, "pick": 317}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 310}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 401}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 316, "wonder": 10}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 308}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 315}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 306, "wonder": 12}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 311, "wonder": 6}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 403}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 303}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 307}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 304}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 400}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 302}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 1}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 301}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 318}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 309}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 313}}]`

	var log []Item

	if err := json.Unmarshal([]byte(raw), &log); err != nil {
		panic("parse log1 panic")
	}

	var out []string

	for _, v := range log {
		m := v.Move.(map[string]interface{})
		idRaw, _ := m["id"]
		id := engine.MoveId(idRaw.(float64))

		switch id {
		case engine.MovePrepare:
			wondersRaw, _ := m["wonders"].([]interface{})
			wonders := []string{}

			for _, widRaw := range wondersRaw {
				wid := widRaw.(float64)

				wonders = append(wonders, wmap[wid])
			}

			tokensRaw, _ := m["tokens"].([]interface{})
			tokens := []string{}

			for _, tidRaw := range tokensRaw {
				tid := tidRaw.(float64)

				tokens = append(tokens, tmap[tid])
			}

			rtokensRaw, _ := m["randomTokens"].([]interface{})
			rtokens := []string{}

			for _, rtidRaw := range rtokensRaw {
				rtid := rtidRaw.(float64)

				rtokens = append(rtokens, tmap[rtid])
			}

			cardsRaw, _ := m["cards"].(map[string]interface{})

			rawcards1 := cardsRaw["1"].([]interface{})
			rawcards2 := cardsRaw["2"].([]interface{})
			rawcards3 := cardsRaw["3"].([]interface{})

			cards1 := []string{}
			cards2 := []string{}
			cards3 := []string{}

			for _, cid1Raw := range rawcards1 {
				cid := cid1Raw.(float64)

				cards1 = append(cards1, cmap[cid])
			}

			for _, cid2Raw := range rawcards2 {
				cid := cid2Raw.(float64)

				cards2 = append(cards2, cmap[cid])
			}

			for _, cid3Raw := range rawcards3 {
				cid := cid3Raw.(float64)

				cards3 = append(cards3, cmap[cid])
			}

			out = append(out, fmt.Sprintf(
				`
PrepareMove{
	move: move{MovePrepare},
	P1:   "%v",
	P2:   "%v",
	Wonders: WonderList{
		%s,
	},
	Tokens: TokenList{
		%s,
	},
	RandomTokens: TokenList{
		%s,
	},
	Cards: map[age]CardList {
		ageI: {
			%s,
		},
		ageII: {
			%s,
		},
		ageIII: {
			%s,
		},
	}`,
				m["p1"],
				m["p2"],
				strings.Join(wonders, ", \n"),
				strings.Join(tokens, ", \n"),
				strings.Join(rtokens, ", \n"),
				strings.Join(cards1, ", \n"),
				strings.Join(cards2, ", \n"),
				strings.Join(cards3, ", \n"),
			),
			)
		case engine.MovePickWonder:
			wid, _ := m["wonder"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickWonder(%s)`,
				wmap[wid],
			))
		}
	}

	fmt.Println(strings.Join(out, ", \n"))
}

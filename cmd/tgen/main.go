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
	raw := `[{"meta": {"actor": ""}, "move": {"id": 1, "p1": "user1", "p2": "user2", "cards": {"1": [116, 119, 122, 114, 120, 109, 112, 106, 101, 100, 121, 104, 103, 102, 117, 115, 105, 113, 118, 111], "2": [215, 208, 209, 207, 203, 201, 216, 217, 220, 212, 213, 222, 218, 210, 202, 214, 205, 200, 211, 204], "3": [305, 302, 309, 314, 310, 307, 317, 306, 403, 400, 311, 304, 300, 301, 319, 318, 315, 316, 405, 308]}, "tokens": [3, 1, 7, 9, 4], "wonders": [6, 12, 3, 13, 9, 14, 7, 10], "randomTokens": [10, 8, 5]}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 12}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 6}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 3}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 13}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 10}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 14}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 7}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 9}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 113}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 111}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 117}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 105}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 104}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 115}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 118}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 102}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 100}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 121}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 103}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 101}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 106}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 120}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 109}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 112}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 122}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 114}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 119}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 116}}, {"meta": {"actor": "user1"}, "move": {"id": 7, "player": "user1"}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 204}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 200}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 202}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 213}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 201}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 211}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 214}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 9}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 205}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 222}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 210}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 218, "wonder": 13}}, {"meta": {"actor": "user1"}, "move": {"id": 10, "card": 215}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 3}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 217}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 1}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 212}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 220}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 203}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 216}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 209}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 207}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 208}}, {"meta": {"actor": "user1"}, "move": {"id": 7, "player": "user1"}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 405, "wonder": 7}}, {"meta": {"actor": "user1"}, "move": {"id": 11, "card": 213}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 318}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 304}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 308}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 315, "wonder": 9}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 300}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 316}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 301}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 400}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 317, "wonder": 14}}, {"meta": {"actor": "user1"}, "move": {"id": 12, "give": 312, "pick": 303}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 307, "wonder": 12}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 309}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 319}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 311}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 403}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 306}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 314}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 310}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 305}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 302}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 7}}]`

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
		case engine.MovePickBoardToken:
			tid, _ := m["token"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickBoardToken(%s)`,
				tmap[tid],
			))
		case engine.MoveConstructCard:
			cid, _ := m["card"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMoveConstructCard(%s)`,
				cmap[cid],
			))
		case engine.MoveConstructWonder:
			cid, _ := m["card"].(float64)
			wid, _ := m["wonder"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMoveConstructWonder(%s, %s)`,
				wmap[wid],
				cmap[cid],
			))
		case engine.MoveDiscardCard:
			cid, _ := m["card"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMoveDiscardCard(%s)`,
				cmap[cid],
			))
		case engine.MoveSelectWhoBeginsTheNextAge:
			p, _ := m["player"].(string)

			out = append(out, fmt.Sprintf(
				`NewMoveSelectWhoBeginsTheNextAge("%s")`,
				p,
			))
		case engine.MoveBurnCard:
			cid, _ := m["card"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMoveBurnCard(%s)`,
				cmap[cid],
			))
		case engine.MovePickRandomToken:
			tid, _ := m["token"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickRandomToken(%s)`,
				tmap[tid],
			))
		case engine.MovePickTopLineCard:
			cid, _ := m["card"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickTopLineCard(%s)`,
				cmap[cid],
			))
		case engine.MovePickDiscardedCard:
			cid, _ := m["card"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickDiscardedCard(%s)`,
				cmap[cid],
			))
		case engine.MovePickReturnedCards:
			pick, _ := m["pick"].(float64)
			give, _ := m["give"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMovePickReturnedCards(%s, %s)`,
				cmap[pick],
				cmap[give],
			))
		case engine.MoveOver:
			loser, _ := m["loser"].(string)
			reason, _ := m["reason"].(float64)

			out = append(out, fmt.Sprintf(
				`NewMoveOver("%s", %s)`,
				loser,
				vmap[reason],
			))
		}
	}

	fmt.Println(strings.Join(out, ", \n"))
}

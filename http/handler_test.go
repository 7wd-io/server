package http

import (
	"7wd.io/di"
	"7wd.io/domain"
	pgsuite "7wd.io/tt/suite/pg"
	"context"
	"fmt"
	swde "github.com/7wd-io/engine"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
	"time"
)

func Test_game(t *testing.T) {
	suite.Run(t, new(gameSuite))
}

type gameSuite struct {
	suite.Suite
	pgs  pgsuite.S
	apis S
	c    *di.C
	//srv *fiber.App
	//svc domain.GameService
}

func (dst *gameSuite) SetupSuite() {
	c := di.MustNew()
	//srv := New()

	gameSvc := domain.NewGameService(
		c.Clock,
		c.Repo.Room,
		c.Repo.Game,
		c.Repo.GameClock,
		c.Repo.User,
		c.Dispatcher,
	)

	dst.pgs.SetupSuite()
	dst.apis.SetupSuite(SuiteOptions{
		Svc:   NewGame(gameSvc),
		Suite: &dst.Suite,
	})

	dst.c = c
	//dst.srv = srv
	//dst.svc = gameSvc

	// создать сервер
	// привязать роуты
	// вызвать fiber.test()

	// mute
}

func (dst *gameSuite) TearDownSuite() {
	dst.pgs.TearDownSuite()
	dst.apis.TearDownSuite()
}

func (dst *gameSuite) SetupTest() {
	dst.pgs.SetupTest(pgsuite.Options{
		Path: path.Join("http", "fixtures", "game"),
	})

	dst.apis.SetupTest()
}

func (dst *gameSuite) TearDownTest() {
	dst.pgs.TearDownTest()
	dst.c.Client.Rds.FlushDB(context.Background())
	dst.apis.TearDownTest()
}

// [{"meta": {"actor": ""}, "move": {"id": 1, "p1": "user1", "p2": "user2", "cards": {"1": [116, 119, 122, 114, 120, 109, 112, 106, 101, 100, 121, 104, 103, 102, 117, 115, 105, 113, 118, 111], "2": [215, 208, 209, 207, 203, 201, 216, 217, 220, 212, 213, 222, 218, 210, 202, 214, 205, 200, 211, 204], "3": [305, 302, 309, 314, 310, 307, 317, 306, 403, 400, 311, 304, 300, 301, 319, 318, 315, 316, 405, 308]}, "tokens": [3, 1, 7, 9, 4], "wonders": [6, 12, 3, 13, 9, 14, 7, 10], "randomTokens": [10, 8, 5]}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 12}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 6}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 3}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 13}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 10}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 14}}, {"meta": {"actor": "user1"}, "move": {"id": 2, "wonder": 7}}, {"meta": {"actor": "user2"}, "move": {"id": 2, "wonder": 9}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 113}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 111}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 117}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 105}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 104}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 115}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 118}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 102}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 100}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 121}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 103}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 101}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 106}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 120}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 109}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 112}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 122}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 114}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 119}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 116}}, {"meta": {"actor": "user1"}, "move": {"id": 7, "player": "user1"}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 204}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 200}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 202}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 213}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 201}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 211}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 214}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 9}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 205}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 222}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 210}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 218, "wonder": 13}}, {"meta": {"actor": "user1"}, "move": {"id": 10, "card": 215}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 3}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 217}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 1}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 212}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 220}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 203}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 216}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 209}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 207}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 208}}, {"meta": {"actor": "user1"}, "move": {"id": 7, "player": "user1"}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 405, "wonder": 7}}, {"meta": {"actor": "user1"}, "move": {"id": 11, "card": 213}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 318}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 304}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 308}}, {"meta": {"actor": "user2"}, "move": {"id": 5, "card": 315, "wonder": 9}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 300}}, {"meta": {"actor": "user2"}, "move": {"id": 6, "card": 316}}, {"meta": {"actor": "user1"}, "move": {"id": 6, "card": 301}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 400}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 317, "wonder": 14}}, {"meta": {"actor": "user1"}, "move": {"id": 12, "give": 312, "pick": 303}}, {"meta": {"actor": "user1"}, "move": {"id": 5, "card": 307, "wonder": 12}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 309}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 319}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 311}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 403}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 306}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 314}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 310}}, {"meta": {"actor": "user2"}, "move": {"id": 4, "card": 305}}, {"meta": {"actor": "user1"}, "move": {"id": 4, "card": 302}}, {"meta": {"actor": "user1"}, "move": {"id": 3, "token": 7}}]
func (dst *gameSuite) Test_Game1() {
	ctx := context.Background()

	user1, err := dst.c.Repo.User.Find(ctx, domain.WithUserNickname("user1"))

	if err != nil {
		dst.FailNow(err.Error())
	}

	if user1 == nil {
		dst.FailNow("game 1: user1 not found")
	}

	user2, err := dst.c.Repo.User.Find(ctx, domain.WithUserNickname("user2"))

	if err != nil {
		dst.FailNow(err.Error())
	}

	if user2 == nil {
		dst.FailNow("game 1: user2 not found")
	}

	now := time.Now()

	o := domain.RoomOptions{
		TimeBank: domain.TimeBankDefault,
	}

	game := &domain.Game{
		HostNickname:  user1.Nickname,
		HostRating:    user1.Rating,
		HostPoints:    domain.Elo(user1.Rating, user2.Rating),
		GuestNickname: user2.Nickname,
		GuestRating:   user2.Rating,
		GuestPoints:   domain.Elo(user2.Rating, user1.Rating),
		Log: domain.GameLog{
			{
				Move: swde.PrepareMove{
					Id: swde.MovePrepare,
					P1: "user1",
					P2: "user2",
					Wonders: swde.WonderList{
						swde.TheHangingGardens,
						swde.TheTempleOfArtemis,
						swde.TheColossus,
						swde.Messe,
						swde.ThePyramids,
						swde.StatueOfLiberty,
						swde.TheMausoleum,
						swde.TheSphinx,
					},
					Tokens: swde.TokenList{
						swde.Economy,
						swde.Agriculture,
						swde.Philosophy,
						swde.Theology,
						swde.Law,
					},
					RandomTokens: swde.TokenList{
						swde.Urbanism,
						swde.Strategy,
						swde.Masonry,
					},
					Cards: map[swde.Age]swde.CardList{
						swde.AgeI: {
							swde.Palisade,
							swde.Theater,
							swde.Tavern,
							swde.Stable,
							swde.Altar,
							swde.Workshop,
							swde.ClayReserve,
							swde.GlassWorks,
							swde.LoggingCamp,
							swde.LumberYard,
							swde.Baths,
							swde.Quarry,
							swde.ClayPit,
							swde.ClayPool,
							swde.Scriptorium,
							swde.Garrison,
							swde.StonePit,
							swde.WoodReserve,
							swde.Pharmacist,
							swde.StoneReserve,
						},
						swde.AgeII: {
							swde.Dispensary,
							swde.CustomHouse,
							swde.CourtHouse,
							swde.Caravansery,
							swde.GlassBlower,
							swde.BrickYard,
							swde.School,
							swde.Laboratory,
							swde.Aqueduct,
							swde.ArcheryRange,
							swde.ParadeGround,
							swde.Brewery,
							swde.Statue,
							swde.HorseBreeders,
							swde.ShelfQuarry,
							swde.Library,
							swde.Walls,
							swde.SawMill,
							swde.Barracks,
							swde.DryingRoom,
						},
						swde.AgeIII: {
							swde.Port,
							swde.Academy,
							swde.Obelisk,
							swde.Observatory,
							swde.Fortifications,
							swde.Palace,
							swde.Senate,
							swde.Armory,
							swde.MagistratesGuild,
							swde.MerchantsGuild,
							swde.SiegeWorkshop,
							swde.ChamberOfCommerce,
							swde.Arsenal,
							swde.Pretorium,
							swde.Arena,
							swde.Lighthouse,
							swde.Gardens,
							swde.Pantheon,
							swde.MoneyLendersGuild,
							swde.TownHall,
						},
					},
				},
			},
		},
		StartedAt: now,
	}

	if err := dst.c.Repo.Game.Save(ctx, game); err != nil {
		dst.FailNow(err.Error())
	}

	gc := &domain.GameClock{
		Id:         game.Id,
		LastMoveAt: now,
		Turn:       domain.Nickname(game.State().Me.Name),
		Values: map[domain.Nickname]domain.TimeBank{
			user1.Nickname: o.Clock(),
			user2.Nickname: o.Clock(),
		},
	}

	if err := dst.c.Repo.GameClock.Save(ctx, gc); err != nil {
		dst.FailNow(err.Error())
	}

	room := &domain.Room{
		Host:        game.HostNickname,
		HostRating:  game.HostRating,
		Guest:       game.GuestNickname,
		GuestRating: game.GuestRating,
		Options:     o,
		GameId:      game.Id,
	}

	if err = dst.c.Repo.Room.Save(ctx, room); err != nil {
		dst.FailNow(err.Error())
	}

	t1, _ := dst.c.TokenFactory.Token(&domain.Passport{
		Id:       user1.Id,
		Nickname: user1.Nickname,
		Rating:   user1.Rating,
		Settings: user1.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(domain.AccessTokenTtl)),
			Subject:   string(user1.Nickname),
		},
	})

	t2, _ := dst.c.TokenFactory.Token(&domain.Passport{
		Id:       user2.Id,
		Nickname: user2.Nickname,
		Rating:   user2.Rating,
		Settings: user2.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(domain.AccessTokenTtl)),
			Subject:   string(user2.Nickname),
		},
	})

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheTempleOfArtemis,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheHangingGardens,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheColossus,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.Messe,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheSphinx,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.StatueOfLiberty,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheMausoleum,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.ThePyramids,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.WoodReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.StoneReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Scriptorium,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.StonePit,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Quarry,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Garrison,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pharmacist,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayPool,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.LumberYard,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Baths,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayPit,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.LoggingCamp,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.GlassWorks,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Altar,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Workshop,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Tavern,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Stable,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Theater,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Palisade,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/select-move").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"user":   "user1",
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.DryingRoom,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.SawMill,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ShelfQuarry,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ParadeGround,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.BrickYard,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Barracks,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Library,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Theology,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Walls,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Brewery,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.HorseBreeders,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.Messe,
			"cardId":   swde.Statue,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-topline-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Dispensary,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Economy,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Laboratory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Agriculture,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ArcheryRange,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Aqueduct,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.GlassBlower,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.School,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.CourtHouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Caravansery,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.CustomHouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/select-move").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"user":   "user1",
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheMausoleum,
			"cardId":   swde.MoneyLendersGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-discarded-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ParadeGround,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Lighthouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ChamberOfCommerce,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.TownHall,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.ThePyramids,
			"cardId":   swde.Gardens,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Arsenal,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pantheon,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pretorium,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.MerchantsGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.StatueOfLiberty,
			"cardId":   swde.Senate,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-returned-cards").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":     game.Id,
			"pickCardId": swde.Study,
			"giveCardId": swde.Circus,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheTempleOfArtemis,
			"cardId":   swde.Palace,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Obelisk,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Arena,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.SiegeWorkshop,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.MagistratesGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Armory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Observatory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Fortifications,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t2).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Port,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Academy,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(t1).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Philosophy,
		}).
		WithAssertStatusOk().
		Send()

	type player struct {
		Name   domain.Nickname `json:"name"`
		Rating domain.Rating   `json:"rating"`
		Points int             `json:"points"`
	}

	type response struct {
		Id       domain.GameId     `json:"id"`
		Host     player            `json:"host"`
		Guest    player            `json:"guest"`
		Clock    *domain.GameClock `json:"clock,omitempty"`
		State    *swde.State       `json:"state"`
		Finished bool              `json:"finished"`
		Log      domain.GameLog    `json:"log"`
	}

	var resp response

	dst.apis.
		GET(fmt.Sprintf("/game/%d", game.Id)).
		WithToken(t1).
		WithAssertStatusOk().
		Send().
		Response(&resp)

	user1Score := swde.Score{
		Civilian:   20,
		Science:    13,
		Commercial: 6,
		Guilds:     0,
		Wonders:    9,
		Tokens:     11,
		Coins:      11,
		Military:   0,
		Total:      70,
	}

	user2Score := swde.Score{
		Civilian:   6,
		Science:    2,
		Commercial: 9,
		Guilds:     10,
		Wonders:    9,
		Tokens:     0,
		Coins:      6,
		Military:   0,
		Total:      42,
	}

	assert.True(dst.T(), resp.Finished)
	assert.Equal(dst.T(), resp.State.Enemy.Score, user1Score)
	assert.Equal(dst.T(), resp.State.Me.Score, user2Score)
}

package http

import (
	"7wd.io/di"
	"7wd.io/domain"
	pgsuite "7wd.io/tt/suite/pg"
	"context"
	swde "github.com/7wd-io/engine"
	"github.com/golang-jwt/jwt/v5"
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

func (dst *gameSuite) Test_Game1() {
	ctx := context.Background()

	// Реквесты:
	// 	- создать игру
	//  - сделать ходы до конца

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

	//now := data.Now()
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

	user10Token, _ := dst.c.TokenFactory.Token(&domain.Passport{
		Id:       user1.Id,
		Nickname: user1.Nickname,
		Rating:   user1.Rating,
		Settings: user1.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(domain.AccessTokenTtl)),
			Subject:   string(user1.Nickname),
		},
	})

	user11Token, _ := dst.c.TokenFactory.Token(&domain.Passport{
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
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheTempleOfArtemis,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheHangingGardens,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheColossus,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.Messe,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheSphinx,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.StatueOfLiberty,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheMausoleum,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.ThePyramids,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.WoodReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.StoneReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Scriptorium,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.StonePit,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Quarry,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Garrison,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pharmacist,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayPool,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.LumberYard,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Baths,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayPit,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.LoggingCamp,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.GlassWorks,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Altar,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Workshop,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ClayReserve,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Tavern,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Stable,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Theater,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Palisade,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/select-move").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"user":   "user1",
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.DryingRoom,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.SawMill,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ShelfQuarry,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ParadeGround,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.BrickYard,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Barracks,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Library,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Theology,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Walls,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Brewery,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.HorseBreeders,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.Messe,
			"cardId":   swde.Statue,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-topline-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Dispensary,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Economy,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Laboratory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Agriculture,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ArcheryRange,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Aqueduct,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.GlassBlower,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.School,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.CourtHouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Caravansery,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.CustomHouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/select-move").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"user":   "user1",
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheMausoleum,
			"cardId":   swde.MoneyLendersGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-discarded-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ParadeGround,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Lighthouse,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.ChamberOfCommerce,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.TownHall,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.ThePyramids,
			"cardId":   swde.Gardens,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Arsenal,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pantheon,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/discard-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Pretorium,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.MerchantsGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.StatueOfLiberty,
			"cardId":   swde.Senate,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-returned-cards").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":     game.Id,
			"pickCardId": swde.Study,
			"giveCardId": swde.Circus,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheTempleOfArtemis,
			"cardId":   swde.Palace,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Obelisk,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Arena,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.SiegeWorkshop,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.MagistratesGuild,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Armory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Observatory,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Fortifications,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Port,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/construct-card").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId": game.Id,
			"cardId": swde.Academy,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/move/pick-board-token").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":  game.Id,
			"tokenId": swde.Philosophy,
		}).
		WithAssertStatusOk().
		Send()
}

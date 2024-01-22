package http

import (
	"7wd.io/di"
	"7wd.io/domain"
	"7wd.io/tt/data"
	"7wd.io/tt/suite/api"
	pgsuite "7wd.io/tt/suite/pg"
	"context"
	swde "github.com/7wd-io/engine"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

func Test_game(t *testing.T) {
	suite.Run(t, new(gameSuite))
}

type gameSuite struct {
	suite.Suite
	pgs  *pgsuite.S
	apis *api.S
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
	dst.apis.SetupSuite(api.SuiteOptions{
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

	user10, err := dst.c.Repo.User.Find(ctx, domain.WithUserNickname("user10"))

	if err != nil {
		dst.FailNow(err.Error())
	}

	if user10 == nil {
		dst.FailNow("game 1: user10 not found")
	}

	user11, err := dst.c.Repo.User.Find(ctx, domain.WithUserNickname("user11"))

	if err != nil {
		dst.FailNow(err.Error())
	}

	if user11 == nil {
		dst.FailNow("game 1: user11 not found")
	}

	now := data.Now()

	o := domain.RoomOptions{
		TimeBank: domain.TimeBankDefault,
	}

	game := &domain.Game{
		HostNickname:  user10.Nickname,
		HostRating:    user10.Rating,
		HostPoints:    domain.Elo(user10.Rating, user11.Rating),
		GuestNickname: user11.Nickname,
		GuestRating:   user11.Rating,
		GuestPoints:   domain.Elo(user11.Rating, user10.Rating),
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
			user10.Nickname: o.Clock(),
			user11.Nickname: o.Clock(),
		},
	}

	if err := dst.c.Repo.GameClock.Save(ctx, gc); err != nil {
		dst.FailNow(err.Error())
	}

	// Send always last called. All asserts before
	//func (dst *Req) Send() {
	//	//res, err := dst.app.Test(dst.toHttpReq(), -1)
	//	res, err := dst.app.Test(dst.toHttpReq())
	//
	//	dst.ss.NoError(err)
	//
	//	if res == nil {
	//		dst.ss.FailNow("response nil")
	//	} else {
	//		defer res.Body.Close()
	//	}
	//
	//	for _, assert := range dst.asserts {
	//		assert(res)
	//	}
	//}

	user10Token, _ := dst.c.TokenFactory.Token(&domain.Passport{
		Id:       user10.Id,
		Nickname: user10.Nickname,
		Rating:   user10.Rating,
		Settings: user10.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(domain.AccessTokenTtl)),
			Subject:   string(user10.Nickname),
		},
	})

	user11Token, _ := dst.c.TokenFactory.Token(&domain.Passport{
		Id:       user11.Id,
		Nickname: user11.Nickname,
		Rating:   user11.Rating,
		Settings: user11.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(domain.AccessTokenTtl)),
			Subject:   string(user11.Nickname),
		},
	})

	dst.apis.
		POST("/game/pick-wonder").
		WithToken(user10Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheTempleOfArtemis,
		}).
		WithAssertStatusOk().
		Send()

	dst.apis.
		POST("/game/pick-wonder").
		WithToken(user11Token).
		WithParams(map[string]interface{}{
			"gameId":   game.Id,
			"wonderId": swde.TheHangingGardens,
		}).
		WithAssertStatusOk().
		Send()
}

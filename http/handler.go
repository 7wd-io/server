package http

import (
	"7wd.io/domain"
	swde "github.com/7wd-io/engine"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewAccount(svc domain.AccountService) Account {
	return Account{
		svc: svc,
	}
}

type Account struct {
	svc domain.AccountService
}

func (dst Account) Bind(app *fiber.App) {
	g := app.Group("/account")

	g.Post("/signup", dst.signup())
	g.Post("/signin", dst.signin())
	g.Post("/logout", dst.logout())
	g.Post("/refresh", dst.refresh())
	g.Put("/settings", dst.updateSettings())
	g.Get("/:nickname", dst.profile())
	g.Get("/:nickname1/vs/:nickname2", dst.profileVersus())
	app.Get("/top", dst.top())
}

func (dst Account) signup() fiber.Handler {
	type request struct {
		Email    domain.Email    `json:"email" validate:"required,email"`
		Password string          `json:"password" validate:"required,min=6,max=32"`
		Nickname domain.Nickname `json:"nickname" validate:"required,nickname"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		return dst.svc.Signup(ctx.Context(), r.Email, r.Password, r.Nickname)
	}
}

func (dst Account) signin() fiber.Handler {
	type request struct {
		Login       string    `json:"login" validate:"required"`
		Password    string    `json:"password" validate:"required"`
		Fingerprint uuid.UUID `json:"fingerprint" validate:"required"`
	}

	type response struct {
		AccessToken  string    `json:"accessToken"`
		RefreshToken uuid.UUID `json:"refreshToken"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		res, err := dst.svc.Signin(ctx.Context(), r.Login, r.Password, r.Fingerprint)

		if err != nil {
			return err
		}

		return ctx.JSON(response{
			AccessToken:  res.Access,
			RefreshToken: res.Refresh,
		})
	}
}

func (dst Account) logout() fiber.Handler {
	type request struct {
		Fingerprint uuid.UUID `json:"fingerprint" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		err := dst.svc.Logout(ctx.Context(), pass, r.Fingerprint)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Account) refresh() fiber.Handler {
	type request struct {
		Fingerprint  uuid.UUID `json:"fingerprint" validate:"required"`
		RefreshToken uuid.UUID `json:"refresh_token" validate:"required"`
	}

	type response struct {
		AccessToken  string    `json:"accessToken"`
		RefreshToken uuid.UUID `json:"refreshToken"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		res, err := dst.svc.Refresh(ctx.Context(), r.RefreshToken, r.Fingerprint)

		if err != nil {
			return err
		}

		return ctx.JSON(response{
			AccessToken:  res.Access,
			RefreshToken: res.Refresh,
		})
	}
}

func (dst Account) updateSettings() fiber.Handler {
	type request struct {
		AnimationSpeed int  `json:"animationSpeed" validate:"min=1,max=5"`
		OpponentJoined bool `json:"opponentJoined"`
		MyTurn         bool `json:"myTurn"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		s := domain.UserSettings{
			Game: domain.GameSettings{
				AnimationSpeed: r.AnimationSpeed,
			},
			Sounds: domain.SoundsSettings{
				OpponentJoined: r.OpponentJoined,
				MyTurn:         r.MyTurn,
			},
		}

		err := dst.svc.UpdateSettings(ctx.Context(), pass, s)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Account) profile() fiber.Handler {
	type response struct {
		Profile domain.UserProfile `json:"profile"`
	}

	return func(ctx *fiber.Ctx) error {
		nickname := domain.Nickname(ctx.Params("nickname"))

		p, err := dst.svc.Profile(ctx.Context(), nickname)

		if err != nil {
			return err
		}

		return ctx.JSON(response{Profile: *p})
	}
}

func (dst Account) profileVersus() fiber.Handler {
	type response struct {
		Profile domain.GamesReport `json:"profile"`
	}

	return func(ctx *fiber.Ctx) error {
		nickname1 := domain.Nickname(ctx.Params("nickname1"))
		nickname2 := domain.Nickname(ctx.Params("nickname2"))

		p, err := dst.svc.ProfileVersus(ctx.Context(), nickname1, nickname2)

		if err != nil {
			return err
		}

		return ctx.JSON(response{Profile: *p})
	}
}

func (dst Account) top() fiber.Handler {
	// @TODO rename
	type response struct {
		Players domain.Top `json:"players"`
	}

	return func(ctx *fiber.Ctx) error {
		top, err := dst.svc.Top(ctx.Context())

		if err != nil {
			return err
		}

		return ctx.JSON(response{Players: top})
	}
}

func NewRoom(svc domain.RoomService) Room {
	return Room{
		svc: svc,
	}
}

type Room struct {
	svc domain.RoomService
}

func (dst Room) Bind(app *fiber.App) {
	g := app.Group("/room")

	g.Get("/", dst.list())
	g.Post("/", dst.create())
	g.Delete("/:id", dst.delete())
	g.Post("/:id/join", dst.join())
	g.Post("/:id/leave", dst.leave())
	g.Post("/:id/kick", dst.kick())
	g.Post("/:id/start", dst.start())
}

func (dst Room) list() fiber.Handler {
	type response struct {
		Items []*domain.Room `json:"items"`
	}

	return func(ctx *fiber.Ctx) error {
		rooms, err := dst.svc.List(ctx.Context())

		if err != nil {
			return err
		}

		return ctx.JSON(response{Items: rooms})
	}
}

func (dst Room) create() fiber.Handler {
	type request struct {
		Fast         bool            `json:"fast,omitempty"`
		MinRating    domain.Rating   `json:"minRating,omitempty" validate:"omitempty,max=2000"`
		Enemy        domain.Nickname `json:"enemy,omitempty" validate:"omitempty,nickname"`
		PromoWonders bool            `json:"promoWonders"`
	}

	type response struct {
		Id domain.RoomId `json:"id"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		room, err := dst.svc.Create(
			ctx.Context(),
			pass,
			domain.RoomOptions{
				Fast:         r.Fast,
				MinRating:    r.MinRating,
				Enemy:        r.Enemy,
				PromoWonders: r.PromoWonders,
			},
		)

		if err != nil {
			return err
		}

		return ctx.JSON(response{Id: room.Id})
	}
}

func (dst Room) delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Delete(ctx.Context(), pass, id)
	}
}

func (dst Room) join() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Join(ctx.Context(), pass, id)
	}
}

func (dst Room) leave() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Leave(ctx.Context(), pass, id)
	}
}

func (dst Room) kick() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Kick(ctx.Context(), pass, id)
	}
}

func (dst Room) start() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Start(ctx.Context(), pass, id)
	}
}

func NewGame(
	game domain.GameService,
	pa domain.PlayAgainService,
) Game {
	return Game{
		game: game,
		pa:   pa,
	}
}

type Game struct {
	game domain.GameService
	pa   domain.PlayAgainService
}

func (dst Game) Bind(app *fiber.App) {
	g := app.Group("/game")

	g.Post("/", dst.createWithBot())
	g.Get("/:id", dst.get())
	g.Get("/units", dst.units())
	g.Get("/:id/state/:index", dst.state())

	g.Post("/move/construct-card", dst.constructCard())
	g.Post("/move/construct-wonder", dst.constructWonder())
	g.Post("/move/discard-card", dst.discardCard())
	g.Post("/move/select-move", dst.selectWhoBeginsTheNextAge())
	g.Post("/move/pick-wonder", dst.pickWonder())
	g.Post("/move/pick-board-token", dst.pickBoardToken())
	g.Post("/move/pick-random-token", dst.pickRandomToken())
	g.Post("/move/burn-card", dst.burnCard())
	g.Post("/move/pick-discarded-card", dst.pickDiscardedCard())
	g.Post("/move/pick-topline-card", dst.pickTopLineCard())
	g.Post("/move/pick-returned-cards", dst.pickReturnedCards())
	g.Post("/move/resign", dst.resign())
	g.Post("/play-again", dst.playAgain())
}

func (dst Game) createWithBot() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pass, _ := usePassport(ctx)

		return dst.game.CreateWithBot(ctx.Context(), pass)
	}
}

func (dst Game) get() fiber.Handler {
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

	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")

		if err != nil {
			return err
		}

		game, err := dst.game.Get(ctx.Context(), domain.GameId(id))

		if err != nil {
			return err
		}

		res := response{
			Id: game.Id,
			Host: player{
				Name:   game.HostNickname,
				Rating: game.HostRating,
				Points: game.HostPoints,
			},
			Guest: player{
				Name:   game.GuestNickname,
				Rating: game.GuestRating,
				Points: game.GuestPoints,
			},
			State:    game.State(),
			Finished: game.IsOver(),
			Log:      game.Log[1:],
		}

		if !game.IsOver() {
			gc, err := dst.game.Clock(ctx.Context(), game.Id)

			if err != nil {
				return err
			}

			res.Clock = gc
		}

		return ctx.JSON(res)
	}
}

func (dst Game) units() fiber.Handler {
	type response struct {
		Cards   swde.CardMap   `json:"cards"`
		Wonders swde.WonderMap `json:"wonders"`
	}

	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(response{
			Cards:   swde.R.Cards,
			Wonders: swde.R.Wonders,
		})
	}
}

func (dst Game) state() fiber.Handler {
	type response struct {
		State swde.State `json:"state"`
	}

	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")

		if err != nil {
			return err
		}

		index, err := ctx.ParamsInt("index")

		if err != nil {
			return err
		}

		state, err := dst.game.State(ctx.Context(), domain.GameId(id), index)

		if err != nil {
			return err
		}

		return ctx.JSON(response{
			State: *state,
		})
	}
}

func (dst Game) constructCard() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Card swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveConstructCard(r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) constructWonder() fiber.Handler {
	type request struct {
		Game   domain.GameId `json:"gameId" validate:"required"`
		Wonder swde.WonderId `json:"wonderId" validate:"required"`
		Card   swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveConstructWonder(r.Wonder, r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) discardCard() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Card swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveDiscardCard(r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) selectWhoBeginsTheNextAge() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		User swde.Nickname `json:"user" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveSelectWhoBeginsTheNextAge(r.User),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickWonder() fiber.Handler {
	type request struct {
		Game   domain.GameId `json:"gameId" validate:"required"`
		Wonder swde.WonderId `json:"wonderId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickWonder(r.Wonder),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickBoardToken() fiber.Handler {
	type request struct {
		Game  domain.GameId `json:"gameId" validate:"required"`
		Token swde.TokenId  `json:"tokenId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickBoardToken(r.Token),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickRandomToken() fiber.Handler {
	type request struct {
		Game  domain.GameId `json:"gameId" validate:"required"`
		Token swde.TokenId  `json:"tokenId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickRandomToken(r.Token),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) burnCard() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Card swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveBurnCard(r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickDiscardedCard() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Card swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickDiscardedCard(r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickTopLineCard() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Card swde.CardId   `json:"cardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickTopLineCard(r.Card),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) pickReturnedCards() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
		Pick swde.CardId   `json:"pickCardId" validate:"required"`
		Give swde.CardId   `json:"giveCardId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMovePickReturnedCards(r.Pick, r.Give),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) resign() fiber.Handler {
	type request struct {
		Game domain.GameId `json:"gameId" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		_, err := dst.game.Move(
			ctx.Context(),
			pass.Nickname,
			r.Game,
			swde.NewMoveOver(swde.Nickname(pass.Nickname), swde.Resign),
		)

		if err != nil {
			return err
		}

		return ctx.JSON(nil)
	}
}

func (dst Game) playAgain() fiber.Handler {
	type request struct {
		Game   domain.GameId `json:"gameId" validate:"required"`
		Answer bool          `json:"answer"`
	}

	return func(ctx *fiber.Ctx) error {
		r := new(request)

		if err := useBodyRequest(ctx, r); err != nil {
			return err
		}

		pass, _ := usePassport(ctx)

		return dst.pa.UpdateById(ctx.Context(), r.Game, pass.Nickname, r.Answer)
	}
}

func NewOnline(svc domain.OnlineService) Online {
	return Online{svc: svc}
}

type Online struct {
	svc domain.OnlineService
}

func (dst Online) Bind(app *fiber.App) {
	g := app.Group("/online")

	g.Get("/", dst.get())
}

func (dst Online) get() fiber.Handler {
	type response struct {
		Data domain.UsersPreview `json:"data"`
	}

	return func(ctx *fiber.Ctx) error {
		users, err := dst.svc.GetAll(ctx.Context())

		if err != nil {
			return err
		}

		return ctx.JSON(response{Data: users})
	}
}

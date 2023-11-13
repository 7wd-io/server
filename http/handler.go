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
	g.Get("/:nickname1/vs/:nickname2", dst.getVersus())
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

func (dst Account) getVersus() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
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
	g.Post("/join/:id", dst.join())
	g.Post("/leave/:id", dst.leave())
	g.Post("/start/:id", dst.start())
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

func (dst Room) start() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(uuid.MustParse(ctx.Params("id")))

		pass, _ := usePassport(ctx)

		return dst.svc.Start(ctx.Context(), pass, id)
	}
}

func NewGame(svc domain.GameService) Game {
	return Game{
		svc: svc,
	}
}

type Game struct {
	svc domain.GameService
}

func (dst Game) Bind(app *fiber.App) {
	g := app.Group("/game")

	g.Get("/:id", dst.get())
	g.Get("/units", dst.units())
	g.Get("/:id/state/:index", dst.state())
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

		game, err := dst.svc.Get(ctx.Context(), domain.GameId(id))

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
			gc, err := dst.svc.Clock(ctx.Context(), game.Id)

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
	type request struct {
		Id    domain.GameId `param:"id" validate:"required,gid"`
		Index int           `param:"index"`
	}

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

		state, err := dst.svc.State(ctx.Context(), domain.GameId(id), index)

		if err != nil {
			return err
		}

		return ctx.JSON(response{
			State: *state,
		})
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

	g.Get("/")
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

package http

import (
	"7wd.io/domain"
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
	g.Post("/settings", dst.updateSettings())
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
		Email       domain.Email `json:"email" validate:"required"`
		Password    string       `json:"password" validate:"required"`
		Fingerprint uuid.UUID    `json:"fingerprint" validate:"required"`
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

		res, err := dst.svc.Signin(ctx.Context(), r.Email, r.Password, r.Fingerprint)

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
	g.Post("/join/:id", dst.join())
	g.Post("/leave/:id", dst.leave())
}

func (dst Room) list() fiber.Handler {
	return nil
}

func (dst Room) create() fiber.Handler {
	type request struct {
		Size int `json:"size" validate:"required,min=2,max=5"`
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
				Size: r.Size,
			},
		)

		if err != nil {
			return err
		}

		return ctx.JSON(response{Id: room.Id})
	}
}

func (dst Room) join() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(ctx.Params("id"))

		pass, _ := usePassport(ctx)

		return dst.svc.Join(ctx.Context(), pass, id)
	}
}

func (dst Room) leave() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := domain.RoomId(ctx.Params("id"))

		pass, _ := usePassport(ctx)

		return dst.svc.Leave(ctx.Context(), pass, id)
	}
}

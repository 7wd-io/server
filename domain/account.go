package domain

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
)

const (
	accessTokenTtl  = 24 * 30 * time.Hour
	refreshTokenTtl = 30 * 24 * time.Hour
	passwordCost    = 10
)

type Email string

type Nickname string

type UserId int32

type User struct {
	Id        UserId
	Email     Email
	Nickname  Nickname
	Password  string
	Settings  UserSettings
	Rating    Rating
	CreatedAt time.Time
}

type UsersPreview map[Nickname]Rating

type UserProfile struct {
	Nickname    Nickname
	Rank        int    `json:"rank"`
	Rating      Rating `json:"rating"`
	GamesReport `json:"games"`
}

type UserSettings struct {
	Game   GameSettings   `json:"game"`
	Sounds SoundsSettings `json:"sounds"`
}

type GameSettings struct {
	AnimationSpeed int `json:"animationSpeed"`
}

type SoundsSettings struct {
	OpponentJoined bool `json:"opponentJoined"`
	MyTurn         bool `json:"myTurn"`
}

type UserOptions struct {
	Tx       Tx
	Id       UserId
	Email    Email
	Nickname Nickname
}

type UserOption func(o *UserOptions)

func WithUserTx(v Tx) UserOption {
	return func(o *UserOptions) {
		o.Tx = v
	}
}

func WithUserId(v UserId) UserOption {
	return func(o *UserOptions) {
		o.Id = v
	}
}

func WithUserEmail(v Email) UserOption {
	return func(o *UserOptions) {
		o.Email = v
	}
}

func WithUserNickname(v Nickname) UserOption {
	return func(o *UserOptions) {
		o.Nickname = v
	}
}

type Passport struct {
	Id       UserId   `json:"id"`
	Nickname Nickname `json:"nickname"`
	Rating   Rating   `json:"rating"`
	jwt.RegisteredClaims
}

type Session struct {
	UserId       UserId
	RefreshToken uuid.UUID
	Fingerprint  uuid.UUID
}

type Token struct {
	Access  string
	Refresh uuid.UUID
}

func NewAccountService(
	userRepo UserRepo,
	pass Password,
	clock Clock,
	tokenf Tokenf,
	uuidf Uuidf,
	sessionRepo SessionRepo,
	analyst Analyst,
) AccountService {
	return AccountService{
		userRepo:    userRepo,
		pass:        pass,
		clock:       clock,
		tokenf:      tokenf,
		uuidf:       uuidf,
		sessionRepo: sessionRepo,
		analyst:     analyst,
	}
}

type AccountService struct {
	userRepo    UserRepo
	pass        Password
	clock       Clock
	tokenf      Tokenf
	uuidf       Uuidf
	sessionRepo SessionRepo
	analyst     Analyst
}

func (dst AccountService) Signup(ctx context.Context, email Email, password string, nickname Nickname) error {
	var err error

	password, err = dst.pass.Hash(password, passwordCost)

	if err != nil {
		return err
	}

	user := &User{
		Email:    email,
		Nickname: nickname,
		Password: password,
		Settings: UserSettings{
			Game: GameSettings{
				AnimationSpeed: 3,
			},
			Sounds: SoundsSettings{
				OpponentJoined: true,
				MyTurn:         false,
			},
		},
		Rating:    1500,
		CreatedAt: dst.clock.Now(),
	}

	return dst.userRepo.Save(ctx, user)
}

func (dst AccountService) Signin(ctx context.Context, login string, pass string, fingerprint uuid.UUID) (*Token, error) {
	var user *User
	var err error

	if strings.Contains(login, "@") {
		user, err = dst.userRepo.Find(ctx, WithUserEmail(Email(login)))
	} else {
		user, err = dst.userRepo.Find(ctx, WithUserNickname(Nickname(login)))
	}

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errCredentialsNotFound
		}

		return nil, err
	}

	if !dst.pass.Check(user.Password, pass) {
		return nil, errCredentialsNotFound
	}

	return dst.token(ctx, user, fingerprint)
}

func (dst AccountService) Logout(ctx context.Context, pass Passport, fingerprint uuid.UUID) error {
	session, err := dst.sessionRepo.Find(ctx, fingerprint)

	if err != nil {
		return err
	}

	if session == nil {
		return nil
	}

	if session.UserId != pass.Id {
		slog.Warn("try logout by user=%d for user=%d", pass.Id, session.UserId)
		return nil
	}

	_, err = dst.sessionRepo.Delete(ctx, fingerprint)

	return err
}

func (dst AccountService) Refresh(ctx context.Context, refreshToken uuid.UUID, fingerprint uuid.UUID) (*Token, error) {
	session, err := dst.sessionRepo.Delete(ctx, fingerprint)

	if err != nil || session == nil {
		return nil, errCredentialsNotFound
	}

	if session.RefreshToken != refreshToken {
		return nil, errCredentialsNotFound
	}

	user, err := dst.userRepo.Find(ctx, WithUserId(session.UserId))

	if err != nil {
		return nil, err
	}

	return dst.token(ctx, user, fingerprint)
}

func (dst AccountService) UpdateSettings(ctx context.Context, pass Passport, s UserSettings) error {
	user, err := dst.userRepo.Find(ctx, WithUserId(pass.Id))

	if err != nil {
		return err
	}

	user.Settings = s

	return dst.userRepo.Update(ctx, user)
}

func (dst AccountService) Profile(ctx context.Context, u Nickname) (*UserProfile, error) {
	gr, err := dst.analyst.GamesReport(ctx, u)

	if err != nil {
		return nil, err
	}

	rank, err := dst.analyst.Rank(ctx, u)

	if err != nil {
		return nil, err
	}

	rating, err := dst.analyst.Rating(ctx, u)

	if err != nil {
		return nil, err
	}

	return &UserProfile{
		Nickname:    u,
		Rank:        rank,
		Rating:      rating,
		GamesReport: *gr,
	}, nil
}

func (dst AccountService) ProfileVersus(ctx context.Context, me Nickname, enemy Nickname) (*GamesReport, error) {
	return dst.analyst.GamesReportVersus(ctx, me, enemy)
}

func (dst AccountService) Top(ctx context.Context) (Top, error) {
	return dst.analyst.Top(ctx)
}

func (dst AccountService) token(ctx context.Context, u *User, fingerprint uuid.UUID) (*Token, error) {
	refresh := dst.uuidf.Uuid()

	access, err := dst.tokenf.Token(&Passport{
		Id:       u.Id,
		Nickname: u.Nickname,
		Rating:   u.Rating,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(dst.clock.Now().Add(accessTokenTtl)),
			Subject:   string(u.Nickname),
		},
	})

	if err != nil {
		return nil, err
	}

	session := &Session{
		UserId:       u.Id,
		RefreshToken: refresh,
		Fingerprint:  fingerprint,
	}

	if err = dst.sessionRepo.Save(ctx, session, refreshTokenTtl); err != nil {
		return nil, err
	}

	return &Token{
		Access:  access,
		Refresh: refresh,
	}, nil
}

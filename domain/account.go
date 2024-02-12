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
	AccessTokenTtl  = 24 * 30 * time.Hour
	RefreshTokenTtl = 30 * 24 * time.Hour
	PasswordCost    = 10
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
	Tx Tx

	Lock bool

	Id    UserId
	IdSet bool

	Email    Email
	EmailSet bool

	Nickname    Nickname
	NicknameSet bool
}

type UserOption func(o *UserOptions)

func WithUserTx(v Tx) UserOption {
	return func(o *UserOptions) {
		o.Tx = v
	}
}

func WithUserLock() UserOption {
	return func(o *UserOptions) {
		o.Lock = true
	}
}

func WithUserId(v UserId) UserOption {
	return func(o *UserOptions) {
		o.Id = v
		o.IdSet = true
	}
}

func WithUserEmail(v Email) UserOption {
	return func(o *UserOptions) {
		o.Email = v
		o.EmailSet = true
	}
}

func WithUserNickname(v Nickname) UserOption {
	return func(o *UserOptions) {
		o.Nickname = v
		o.NicknameSet = true
	}
}

type Passport struct {
	Id       UserId       `json:"id"`
	Nickname Nickname     `json:"nickname"`
	Rating   Rating       `json:"rating"`
	Settings UserSettings `json:"settings"`
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
	tx Txer,
) AccountService {
	return AccountService{
		userRepo:    userRepo,
		pass:        pass,
		clock:       clock,
		tokenf:      tokenf,
		uuidf:       uuidf,
		sessionRepo: sessionRepo,
		analyst:     analyst,
		tx:          tx,
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
	tx          Txer
}

func (dst AccountService) Signup(ctx context.Context, email Email, password string, nickname Nickname) error {
	var err error

	password, err = dst.pass.Hash(password, PasswordCost)

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
		Rating:    DefaultElo,
		CreatedAt: dst.clock.Now(),
	}

	if err = dst.userRepo.Save(ctx, user); err != nil {
		return err
	}

	return dst.analyst.UpdateRatings(ctx, user)
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
		slog.Warn(
			"AccountService.Logout session.UserId != pass.Id",
			slog.Int("from", int(pass.Id)),
			slog.Int("for", int(session.UserId)),
		)
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
	tx, err := dst.tx.Tx(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			slog.Error("AccountService.UpdateSettings: tx.Rollback", "err", err)
		}
	}()

	user, err := dst.userRepo.Find(
		ctx,
		WithUserId(pass.Id),
		WithUserTx(tx),
		WithUserLock(),
	)

	if err != nil {
		return err
	}

	user.Settings = s

	if err = dst.userRepo.Update(ctx, user, WithUserTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
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

func (dst AccountService) OnGameOver(ctx context.Context, payload interface{}) error {
	p, ok := payload.(GameOverPayload)

	if !ok {
		return errors.New("func (dst AccountService) OnGameOver !ok := payload.(GameOverPayload)")
	}

	var err error

	wBot := p.Result.Winner == BotNickname
	lBot := p.Result.Loser == BotNickname

	botGame := wBot || lBot

	if !botGame {
		err = dst.updateRating(ctx, p.Result.Winner, p.Result.Points)

		if err != nil {
			return err
		}

		err = dst.updateRating(ctx, p.Result.Loser, -p.Result.Points)

		if err != nil {
			return err
		}
	} else {
		bot := p.Result.Winner
		points := p.Result.Points

		if lBot {
			bot = p.Result.Loser
			points = -points
		}

		err = dst.updateRating(ctx, bot, points)

		if err != nil {
			return err
		}
	}

	return nil
}

func (dst AccountService) token(ctx context.Context, u *User, fingerprint uuid.UUID) (*Token, error) {
	refresh := dst.uuidf.Uuid()

	access, err := dst.tokenf.Token(&Passport{
		Id:       u.Id,
		Nickname: u.Nickname,
		Rating:   u.Rating,
		Settings: u.Settings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(dst.clock.Now().Add(AccessTokenTtl)),
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

	if err = dst.sessionRepo.Save(ctx, session, RefreshTokenTtl); err != nil {
		return nil, err
	}

	return &Token{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (dst AccountService) updateRating(ctx context.Context, u Nickname, points int) error {
	var err error

	tx, err := dst.tx.Tx(ctx)

	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("AccountService.updateRating: tx.Rollback", "err", err)
		}
	}()

	user, _ := dst.userRepo.Find(
		ctx,
		WithUserNickname(u),
		WithUserTx(tx),
		WithUserLock(),
	)
	user.Rating += Rating(points)

	if err = dst.userRepo.Update(ctx, user, WithUserTx(tx)); err != nil {
		return err
	}

	if err = dst.analyst.UpdateRatings(ctx, user); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

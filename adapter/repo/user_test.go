package repo

import (
	"7wd.io/domain"
	"7wd.io/tt/data"
	"7wd.io/tt/pg"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"path"
	"testing"
)

func Test_User(t *testing.T) {
	suite.Run(t, new(userSuite))
}

type userSuite struct {
	suite.Suite
	pg pg.Suite
	r  UserRepo
}

func (dst *userSuite) SetupSuite() {
	dst.pg.SetupSuite()

	dst.r = NewUser(dst.pg.C)
}

func (dst *userSuite) TearDownSuite() {
	dst.pg.TearDownSuite()
}

func (dst *userSuite) SetupTest() {
	dst.pg.SetupTest(pg.TestOptions{
		Path: path.Join("adapter", "repo", "fixtures", "user"),
	})
}

func (dst *userSuite) TearDownTest() {
	dst.pg.TearDownTest()
}

func (dst *userSuite) Test_New() {
	assert.NotNil(dst.T(), NewUser(dst.pg.C))
}

func (dst *userSuite) Test_Find() {
	createdAt := data.Now()

	expected := &domain.User{
		Id:       10,
		Rating:   1500,
		Nickname: "User10",
		Email:    "user10@gmail.com",
		Password: "12345678",
		Settings: domain.UserSettings{
			Game: domain.GameSettings{
				AnimationSpeed: 3,
			},
			Sounds: domain.SoundsSettings{
				MyTurn:         false,
				OpponentJoined: true,
			},
		},
		CreatedAt: createdAt,
	}

	actual, err := dst.r.Find(context.Background(), domain.WithUserId(10))

	assert.NoError(dst.T(), err)
	assert.Equal(dst.T(), expected, actual)
}

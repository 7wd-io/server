package http

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test_game(t *testing.T) {
	suite.Run(t, new(gameSuite))
}

type gameSuite struct {
	suite.Suite
}

func (dst *gameSuite) SetupSuite() {
	// mute
}

func (dst *gameSuite) TearDownSuite() {
	// mute
}

func (dst *gameSuite) SetupTest() {
	// mute
}

func (dst *gameSuite) TearDownTest() {
	// mute
}

func (dst *gameSuite) Test_Game1() {

}

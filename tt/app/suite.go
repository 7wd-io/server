package app

import (
	"7wd.io/http"
	"7wd.io/tt/suite/pg"
	"github.com/stretchr/testify/suite"
	"path"
)

type Suite struct {
	suite.Suite
	pg          pg.S
	API         http.S
	fixturesDir string
}

func (dst *Suite) SetupSuite(o SuiteOptions) {
	dst.fixturesDir = o.FixturesDir

	dst.API.SetupSuite(http.SuiteOptions{
		Svc:   o.Svc,
		Suite: &dst.Suite,
	})

	dst.pg.SetupSuite()
}

func (dst *Suite) TearDownSuite() {
	dst.API.TearDownSuite()
	dst.pg.TearDownSuite()
}

func (dst *Suite) SetupTest(o TestOptions) {
	dst.API.SetupTest()

	pgOptions := pg.Options{}

	if o.Fixtures != "" {
		pgOptions.Path = path.Join(dst.fixturesDir, o.Fixtures)
	}

	dst.pg.SetupTest(pgOptions)
}

func (dst *Suite) TearDownTest() {
	dst.API.TearDownTest()
	dst.pg.TearDownTest()
}

type SuiteOptions struct {
	Svc         Binder
	BaseUrl     string
	FixturesDir string
}

type TestOptions struct {
	Fixtures string
}

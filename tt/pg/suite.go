package pg

import (
	"7wd.io/config"
	"7wd.io/fixtures"
	"7wd.io/infra/pg"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"path"
	"strings"
)

type Suite struct {
	C *pgxpool.Pool
}

func (dst *Suite) SetupSuite() {
	dst.C = pg.MustNew(context.Background())
}

func (dst *Suite) TearDownSuite() {
	dst.C.Close()
}

func (dst *Suite) SetupTest(o TestOptions) {
	if o.Path != "" {
		dir := path.Join(config.C.Path, o.Path)

		if err := fixtures.MustNew(dir).Load(); err != nil {
			log.Fatal(errRoot(fmt.Errorf("load: %w", err)))
		}
	}
}

func (dst *Suite) TearDownTest() {
	dst.clear()
}

type TestOptions struct {
	Path string
}

func (dst *Suite) clear() {
	tables := []string{
		`"user"`,
		`"game"`,
	}

	_, err := dst.C.Exec(
		context.Background(),
		fmt.Sprintf("TRUNCATE %s RESTART IDENTITY", strings.Join(tables, ", ")),
	)

	if err != nil {
		log.Fatal(errRoot(fmt.Errorf("clear: %w", err)))
	}
}

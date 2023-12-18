package fixtures

import (
	"7wd.io/config"
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func MustNew(dir string) *testfixtures.Loader {
	if !config.C.IsTest() {
		log.Fatalln("only for test environment")
	}

	db, err := sql.Open(
		"pgx",
		config.C.PgDsn(),
	)

	if err != nil {
		log.Fatalln(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dir),
		//testfixtures.UseAlterConstraint(),
		//testfixtures.UseDropConstraint(),
		testfixtures.ResetSequencesTo(1),
	)

	if err != nil {
		log.Fatalln(err)
	}

	return fixtures
}

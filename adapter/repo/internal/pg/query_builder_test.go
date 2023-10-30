package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test_QB(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(qbSuite))
}

type qbSuite struct {
	suite.Suite
	qb QB
}

func (dst *qbSuite) SetupTest() {
	dst.qb = QB{
		TableName: "tname",
		Columns: []string{
			"id",
			"col1",
			"col2",
		},
	}
}

func (dst *qbSuite) Test_GetSelectQuery() {
	expected := `SELECT id,col1,col2 FROM tname WHERE TRUE`

	assert.Equal(dst.T(), expected, dst.qb.Select(" WHERE TRUE"))
}

func (dst *qbSuite) Test_SelectWhere_when_empty() {
	expected := `SELECT id,col1,col2 FROM tname`

	assert.Equal(dst.T(), expected, dst.qb.SelectWhere(Where{}))
}

func (dst *qbSuite) Test_SelectWhere_when_single() {
	expected := `SELECT id,col1,col2 FROM tname WHERE id=$1`

	assert.Equal(dst.T(), expected, dst.qb.SelectWhere(Where{F{Expr: "id", Value: 1}}))
}

func (dst *qbSuite) Test_SelectWhereEqual_when_multi() {
	expected := `SELECT id,col1,col2 FROM tname WHERE col1=$1 AND col2=$2`

	assert.Equal(dst.T(), expected, dst.qb.SelectWhere(Where{
		F{Expr: "col1", Value: 1},
		F{Expr: "col2", Value: 2},
	}))
}

func (dst *qbSuite) Test_Insert() {
	expected := `INSERT INTO tname (col1,col2) VALUES($1,$2) RETURNING id`

	assert.Equal(dst.T(), expected, dst.qb.Insert())
}

func (dst *qbSuite) Test_Update() {
	expected := `UPDATE tname SET col1=$2,col2=$3 WHERE id=$1`

	assert.Equal(dst.T(), expected, dst.qb.Update())
}

func (dst *qbSuite) Test_Delete() {
	expected := `DELETE FROM tname WHERE id=$1`

	assert.Equal(dst.T(), expected, dst.qb.Delete())
}

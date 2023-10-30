package pg

import (
	"fmt"
	"strings"
)

type QB struct {
	TableName string
	Columns   []string
}

func (dst QB) Select(afterQuery string) string {
	return fmt.Sprintf(
		`SELECT %s FROM %s%s`,
		dst.columns(),
		dst.TableName,
		afterQuery,
	)
}

func (dst QB) SelectWhere(c Where) string {
	where := c.String()

	if where == "" {
		return dst.Select("")
	}

	return dst.Select(fmt.Sprintf(` WHERE %s`, where))
}

func (dst QB) Insert() string {
	cols := dst.Columns[1:]
	l := len(cols)

	var p []string

	for i := 1; i <= l; i++ {
		p = append(p, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(%s) RETURNING id`,
		dst.TableName,
		strings.Join(cols, ","),
		strings.Join(p, ","),
	)
}

func (dst QB) Update() string {
	var q []string

	pos := 2
	for _, c := range dst.Columns[1:] {
		q = append(q, fmt.Sprintf("%s=$%d", c, pos))
		pos++
	}

	return fmt.Sprintf(
		`UPDATE %s SET %s WHERE id=$1`,
		dst.TableName,
		strings.Join(q, ","),
	)
}

func (dst QB) Delete() string {
	return fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, dst.TableName)
}

func (dst QB) columns() string {
	return strings.Join(dst.Columns, ",")
}

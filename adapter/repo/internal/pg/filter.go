package pg

import (
	"fmt"
	"strings"
)

type Where []Filter

func (dst Where) String() string {
	var items []string

	var n int

	for _, item := range dst {
		items = append(items, item.String(&n))
	}

	return strings.Join(items, " AND ")
}

func (dst Where) Values() []interface{} {
	var values []interface{}

	for _, value := range dst {
		for _, v := range value.Values() {
			if v != nil {
				values = append(values, v)
			}
		}
	}

	return values
}

type F struct {
	// column name or raw sql
	Expr string
	// def =
	Op    string
	Value interface{}
	// def $%d
	ValueT string
}

func (dst F) String(n *int) string {
	var valueT, op string

	if dst.Value == nil {
		return dst.Expr
	}

	*n++

	valueT = `$%d`

	if dst.ValueT != "" {
		valueT = dst.ValueT
	}

	op = "="

	if dst.Op != "" {
		op = dst.Op
	}

	return fmt.Sprintf(`%s%s`+valueT, dst.Expr, op, *n)
}

func (dst F) Values() []interface{} {
	return []interface{}{dst.Value}
}

type OR struct {
	Items []F
}

func (dst OR) String(n *int) string {
	items := make([]string, len(dst.Items))

	for k, v := range dst.Items {
		items[k] = v.String(n)
	}

	return "(" + strings.Join(items, " OR ") + ")"
}

func (dst OR) Values() []interface{} {
	var values []interface{}

	for _, v := range dst.Items {
		if v.Value != nil {
			values = append(values, v.Value)
		}
	}

	return values
}

type Filter interface {
	String(n *int) string
	Values() []interface{}
}

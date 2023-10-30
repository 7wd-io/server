package pg

import (
	"reflect"
	"testing"
)

func TestF_String(t *testing.T) {
	tests := []struct {
		name string
		dst  Where
		want string
	}{
		{
			name: "Smoke",
			dst: Where{
				F{
					Expr:  "col1",
					Value: "some string",
				},
				F{
					Expr:  "col2",
					Value: nil,
				},
				F{
					Expr:  "col3 IS NULL",
					Value: nil,
				},
				F{
					Expr:  "col4",
					Value: 1,
				},
				F{
					Expr:   `"col5"`,
					Value:  []int{1, 2, 3},
					ValueT: "ANY($%d)",
				},
				OR{
					Items: []F{
						{
							Expr:  "col6",
							Value: 1,
						},
						{
							Expr:  "col7",
							Value: 1,
						},
					},
				},
				F{
					Expr:   "col8",
					Op:     "@>",
					Value:  10,
					ValueT: "ARRAY [$%d]",
				},
			},
			want: `col1=$1 AND col2 AND col3 IS NULL AND col4=$2 AND "col5"=ANY($3) AND (col6=$4 OR col7=$5) AND col8@>ARRAY [$6]`,
		},
		{
			name: "Empty",
			dst:  Where{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dst.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestF_Values(t *testing.T) {
	tests := []struct {
		name string
		dst  Where
		want []interface{}
	}{
		{
			name: "Smoke",
			dst: Where{
				F{
					Expr:  "col1",
					Value: "some string",
				},
				F{
					Expr: "col2",
				},
				F{
					Expr: "col3 IS NULL",
				},
				OR{
					Items: []F{
						{
							Expr:  "col4",
							Value: "OR1",
						},
						{
							Expr: "col5",
						},
						{
							Expr:  "col6",
							Value: "OR2",
						},
					},
				},
				F{
					Expr:  "col7",
					Value: 1,
				},
			},
			want: []interface{}{"some string", "OR1", "OR2", 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dst.Values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

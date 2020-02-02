package store

import "context"

type List interface {
	Columns() []string
}

type Option func(q Query) Query

type Query struct {
	ctx        context.Context
	query      string
	searchText string
	parameters map[string]interface{}
	limit      int
	offset     int
}

func (q *Query) Exec(i interface{}) {

}

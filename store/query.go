package store

import "context"

type Option func(q Query) Query

type Query struct {
	ctx    context.Context
	filter string
	limit  int
	offset int
}

func Select(opts ...Option) Query {
	q := Query{}

	for _, opt := range opts {
		q = opt(q)
	}

	return q
}

func Offset(offset int) Option {
	return func(q Query) Query {
		q.offset = offset
		return q
	}
}

func Limit(limit int) Option {
	return func(q Query) Query {
		q.limit = limit
		return q
	}
}

func Filter(filter string) Option {
	return func(q Query) Query {
		q.filter = filter
		return q
	}
}

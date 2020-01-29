package model

import "github.com/webitel/engine/model"

const (
	PAGE_DEFAULT     = 0
	PER_PAGE_DEFAULT = 40
	PER_PAGE_MAXIMUM = 1000
)

type ListRequest struct {
	Q       *string
	page    int
	perPage int
	data    []interface{}
}

func (l *ListRequest) Split(data []interface{}) []interface{} {
	if len(data) > l.perPage {
		return data[:l.perPage]
	}
	return data
}

func (l *ListRequest) SetPage(page int) *ListRequest {
	l.page = page
	l.valid()
	return l
}

func (l *ListRequest) SetQ(q string) *ListRequest {
	if q == "" {
		l.Q = nil
	} else {
		l.Q = model.NewString(q)
	}
	return l
}

func (l *ListRequest) SetPerPage(perPage int) *ListRequest {
	l.perPage = perPage
	l.valid()
	return l
}

func (l *ListRequest) GetLimit() int {
	l.valid()
	return l.perPage + 1 //FIXME for next page...
}

func (l *ListRequest) GetOffset() int {
	l.valid()
	return l.perPage * l.page
}

func (l *ListRequest) valid() {
	if l.page < 0 {
		l.page = PAGE_DEFAULT
	}

	if l.perPage < 1 || l.perPage > PER_PAGE_MAXIMUM {
		l.perPage = PER_PAGE_DEFAULT
	}
}

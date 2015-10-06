package pop

import (
	"encoding/json"
	"strconv"
)

type Paginator struct {
	// Current page you're on
	Page int `json:"page"`
	// Number of results you want per page
	PerPage int `json:"per_page"`
	// Page * PerPage (ex: 2 * 20, Offset == 40)
	Offset int `json:"offset"`
	// Total potential records matching the query
	TotalEntriesSize int `json:"total_entries_size"`
	// Total records returns, will be <= PerPage
	CurrentEntriesSize int `json:"current_entries_size"`
	// Total pages
	TotalPages int `json:"total_pages"`
}

func (p Paginator) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewPaginator(page int, per_page int) *Paginator {
	p := &Paginator{Page: page, PerPage: per_page}
	p.Offset = (p.Page - 1) * p.PerPage
	return p
}

type PaginationParams interface {
	Get(key string) string
}

func NewPaginatorFromParams(params PaginationParams) *Paginator {
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	per_page := params.Get("per_page")
	if per_page == "" {
		per_page = "20"
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	pp, err := strconv.Atoi(per_page)
	if err != nil {
		pp = 20
	}
	return NewPaginator(p, pp)
}

func (c *Connection) Paginate(page int, per_page int) *Query {
	return Q(c).Paginate(page, per_page)
}

func (q *Query) Paginate(page int, per_page int) *Query {
	q.Paginator = NewPaginator(page, per_page)
	return q
}

func (c *Connection) PaginateFromParams(params PaginationParams) *Query {
	return Q(c).PaginateFromParams(params)
}

func (q *Query) PaginateFromParams(params PaginationParams) *Query {
	q.Paginator = NewPaginatorFromParams(params)
	return q
}

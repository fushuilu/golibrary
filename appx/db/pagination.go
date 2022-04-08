package db

import (
	"github.com/fushuilu/golibrary/appx/datax"
	"time"

	"github.com/fushuilu/golibrary"
)

// 分页相关
const MaxPageSize = 30
const DefaultPageSize = 15

var DefaultPagination = Pagination{
	Before: 0, After: 0, PageSize: DefaultPageSize, PageIndex: 0,
}

type Pagination struct {
	Before   int    `json:"before" query:"before"`       // 时间（优先）查询比 Before 小的记录, db.xxx < pag.before
	After    int    `json:"after" query:"after"`         // 时间（优先）查询比 After  大的记录, pag.after < db.xxx
	BeforeAt string `json:"before_at" query:"before_at"` // 时间
	AfterAt  string `json:"after_at" query:"after_at"`   // 时间
	// 分页
	PageSize  int         `json:"page_size" query:"page_size"`   // 每页显示记录数
	PageIndex int         `json:"page_index" query:"page_index"` // 当前页，默认为0
	PageNext  interface{} `json:"page_next" query:"page_next"`   // 下一页
	LimitSize int         `json:"limit" query:"limit"`           // 每页显示记录数 PageSize
	Page      int         `json:"page" query:"page"`             // 当前页 PageIndex，默认为0
}

func (p *Pagination) Offset() int {
	if p.Page > 0 {
		p.PageIndex = p.Page
	}
	return p.Limit() * p.PageIndex
}
func (p *Pagination) Limit() int {
	if p.LimitSize > 0 {
		return p.LimitSize
	}
	return p.PageSize
}
func (p *Pagination) Start() int {
	return p.Offset()
}

func (p *Pagination) IsFirst() bool {
	if p.Page > 0 {
		p.PageIndex = p.Page
	}
	return p.PageIndex == 0
}

func (p *Pagination) BeforeTimeOf() (time.Time, bool) {
	if p.Before > 0 {
		p.BeforeAt = golibrary.AnyToString(p.Before)
	}
	if p.BeforeAt != "" {
		if date, err := golibrary.DateTimeParse(p.BeforeAt); err == nil {
			return date.AddDate(0, 0, 1), true
		}
	}
	return datax.EmptyTime, false
}

func (p *Pagination) AfterTimeOf() (time.Time, bool) {
	if p.After > 0 {
		p.AfterAt = golibrary.AnyToString(p.After)
	}
	if p.AfterAt != "" {
		if date, err := golibrary.DateTimeParse(p.AfterAt); err == nil {
			return date, true
		}
	}
	return datax.EmptyTime, false
}

func (p *Pagination) Where(where *Where, colName string) {
	if after, has := p.AfterTimeOf(); has {
		where.TimeLT(colName, after)
	}
	if before, has := p.BeforeTimeOf(); has {
		where.TimeGTE(colName, before)
	}
}

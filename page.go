package sqly

import (
	"strconv"
	"strings"
)

//type Page struct {
//	Page    int64
//	Size    int64
//	OrderBy string
//}
//
//func (p *Page) GetSize() int64 {
//	if p.Size < 1 {
//		p.Size = 10
//	}
//	if p.Size > 10000 {
//		p.Size = 10000
//	}
//	return p.Size
//}
//func (p *Page) GetPage() int64 {
//	if p.Page < 1 {
//		p.Page = 1
//	}
//	return p.Page
//}
//
//func (p *Page) GetLimit() string {
//	return " limit " + strconv.FormatInt(p.GetOffset(), 10) + " , " + strconv.FormatInt(p.GetSize(), 10)
//
//}
//func (p *Page) GetOffset() int64 {
//	return (p.GetPage() - 1) * p.GetSize()
//}

type Page interface {
	GetSize() int64
	GetPage() int64
	GetOrderBy() string
}

func sqlFormatPage(bindType int, sql string, page Page) string {
	if page.GetOrderBy() != "" {
		sql += " order by " + page.GetOrderBy()
	}
	switch bindType {
	// oracle only supports named type bind vars even for positional
	case NAMED:

	case QUESTION, UNKNOWN:
		return sql + " limit " + strconv.FormatInt((page.GetPage()-1)*page.GetSize(), 10) + " , " + strconv.FormatInt(page.GetSize(), 10)
	case DOLLAR:

	case AT:

	}
	return sql
}

func sqlFormatCount(sql string) string {
	str := strings.ReplaceAll(sql, "\n", " ")
	str = strings.ReplaceAll(str, "\t", " ")
	index := strings.Index(str, "select ")
	count := 1

	for count != 0 {
		i, i2 := inquireSelectOrFrom(str, index)
		index += i
		count += i2
	}
	return "select count(*) " + sql[index:]
}

func inquireSelectOrFrom(str string, index int) (int, int) {
	str = strings.ToLower(str)
	selectStr := " select "
	fromStr := " from "

	si := strings.Index(str[index+1:], selectStr)
	fi := strings.Index(str[index+1:], fromStr)

	if si > fi || si == -1 {
		return fi + 1, -1
	} else {
		return si + 1, 1
	}

}

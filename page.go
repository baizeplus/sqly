package sqly

import (
	"strconv"
	"strings"
)

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
	case QUESTION, UNKNOWN, DOLLAR:
		return sql + " limit " + strconv.FormatInt(page.GetSize(), 10) + " offset " + strconv.FormatInt((page.GetPage()-1)*page.GetSize(), 10)
		// oracle only supports named type bind vars even for positional
	case NAMED, AT:
		return sql + " OFFSET " + strconv.FormatInt((page.GetPage()-1)*page.GetSize(), 10) + " ROWS FETCH NEXT " + strconv.FormatInt(page.GetSize(), 10) + " ROWS ONLY"
	}

	return sql
}

func sqlFormatCount(sql string) string {
	str := strings.ReplaceAll(strings.ReplaceAll(sql, "\n", " "), "\t", " ")
	lowerStr := strings.ToLower(str)
	index := strings.Index(lowerStr, "select ")
	count := 1

	for count != 0 {
		i, i2 := inquireSelectOrFrom(lowerStr, index)
		index += i
		count += i2
	}
	return "select count(*) " + sql[index:]
}

func inquireSelectOrFrom(lowerStr string, index int) (int, int) {
	si := strings.Index(lowerStr[index+1:], " select ")
	fi := strings.Index(lowerStr[index+1:], " from ")

	if si > fi || si == -1 {
		return fi + 1, -1
	}
	return si + 1, 1
}

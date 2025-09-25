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
	str := strings.ReplaceAll(sql, "\n", " ")
	str = strings.ReplaceAll(str, "\t", " ")
	lowerStr := strings.ToLower(str)
	
	// Find the main FROM clause
	selectIndex := strings.Index(lowerStr, "select")
	if selectIndex == -1 {
		return sql
	}
	
	parenCount := 0
	for i := selectIndex + 6; i < len(str)-4; i++ {
		if str[i] == '(' {
			parenCount++
		} else if str[i] == ')' {
			parenCount--
		} else if parenCount == 0 && lowerStr[i:i+4] == "from" {
			// Check if it's a complete word (not part of another word)
			if (i == 0 || !isAlphaNum(str[i-1])) && (i+4 >= len(str) || !isAlphaNum(str[i+4])) {
				return "select count(*) " + str[i:]
			}
		}
	}
	return sql
}


func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
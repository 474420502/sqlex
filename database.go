package sqlex

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// TableName 用于tag声明结构的表名
type TableName interface {
}

// DB 继承重载sql.DB
type DB struct {
	*sql.DB
}

// Open 返回sqlex.DB.  sql.DB 被 sqlex.DB的覆盖.
func Open(driverName string, dataSourceName string) (*DB, error) {
	db := &DB{}

	return db, nil
}

// OpenDB 返回sqlex.DB sql.DB 被 sqlex.DB的覆盖.
func OpenDB(c driver.Connector) *DB {
	db := &DB{}

	return db
}

func init() {

	db, err := sql.Open("mysql", "root:Nono-databoard@tcp(sg-board1.livenono.com:3306)/databoard?parseTime=true&loc=Local&charset=utf8&collation=utf8_unicode_ci&timeout=5s")
	if err != nil {
		panic(err)
	}
	db2 := &DB{DB: db}
	row := db2.QueryRow("select count(1) from count_live_anchors")
	var c int64

	row.Scan(&c)
	log.Println(c)
}

type table struct {
	name    string
	index   int
	columns []string
}

type values struct {
}

func getField(query []rune, end rune, i *int) string {
	var field []rune
	*i++ // ( +1 => content
	for ; *i < len(query)-1; *i++ {
		c := query[*i]
		if c == ' ' {
			continue
		}
		if c == end {
			*i++ // ) => next
			break
		}
		field = append(field, c)
	}
	return string(field)
}

func getTable(query []rune, i *int) *table {
	*i += 2 // => ()[]<>
	tb := &table{}
EXTRACTOR_TABLE:
	for {

		switch cquto := query[*i]; cquto {
		case '<':
			tb.name = getField(query, '>', i)
		case '(':
			columnsStr := getField(query, ')', i)
			tb.columns = strings.Split(columnsStr, ",")
		case '[':
			index, err := strconv.Atoi(getField(query, ']', i))
			if err != nil {
				panic(err)
			}
			tb.index = index
		default:
			break EXTRACTOR_TABLE
		}

	}

	return tb
}

func check(q *string) string {
	// prefix := "?t"
	query := []rune(*q)
	var tables []*table
	for i := 0; i < len(query)-1; i++ {
		switch c := rune(query[i]); c {
		case '\\':
			i++
			continue
		case '?':
			inputType := query[i+1]
			switch inputType {
			case 't': // 提取表结构
				tables = append(tables, getTable(query, &i))
			case 'v': // 提取值结构

			}
		}

	}
	// name := "?t<>[]()"
	return ""
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

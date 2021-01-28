package usql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const SQLCOL = "sqlcol"

type DBItem interface {
	TableName() string
	ScanRow(row RowScanner) error
}

// The query starts from the WHERE clause. Returns `NotFound` if no such an record in DB.
func SelectItem(stub Querier, whereQuery string, item DBItem, args ...interface{}) error {
	row := stub.QueryRowContext(newQueryCtx(), selectQuery(item, whereQuery), args...)
	err := item.ScanRow(row)
	if err == sql.ErrNoRows {
		return NotFound
	}
	if err != nil {
		return err
	}
	return nil
}

// The query starts from the WHERE clause.
func SelectItems(stub Querier, whereQuery string, item DBItem, args ...interface{}) ([]DBItem, error) {
	rows, err := stub.QueryContext(newQueryCtx(), selectQuery(item, whereQuery), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DBItem
	for rows.Next() {
		i := newItem(item)
		if err := i.ScanRow(rows); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func newItem(item DBItem) DBItem {
	return reflect.New(reflect.ValueOf(item).Elem().Type()).Interface().(DBItem)
}

// extract column names from struct tag `sqlcol`
// DBItem is expected to be pointer type
func colStr(item DBItem) string {
	t := reflect.TypeOf(item).Elem()
	var cols []string
	for i := 0; i < t.NumField(); i++ {
		col := t.Field(i).Tag.Get(SQLCOL)
		if col != "" {
			cols = append(cols, col)
		}
	}

	if len(cols) != 0 {
		return strings.Join(cols, ",")
	}
	return ""
}

func selectQuery(item DBItem, whereQuery string) string {
	return fmt.Sprintf("SELECT %s FROM %s %s", colStr(item), item.TableName(), whereQuery)
}

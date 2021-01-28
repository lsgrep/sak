package usql

import (
	"database/sql"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

type TestItem struct {
	OpIndex    int64     `sqlcol:"op_index"`
	OpName     string    `sqlcol:"op_name"`
	LogicTime  time.Time `sqlcol:"logic_time"`
	EventIndex int64     `sqlcol:"event_index"`
	InputData  string    `sqlcol:"input_data"`

	TxHash      sql.NullString `sqlcol:"tx_hash"`
	SendTime    time.Time      `sqlcol:"send_time"`
	Status      string         `sqlcol:"status"`
	BlockNumber sql.NullInt64  `sqlcol:"block_number"`
}

func (ti *TestItem) TableName() string {
	return "testItems"
}

func (ti *TestItem) ScanRow(row RowScanner) error {
	return row.Scan(
		&ti.OpIndex, &ti.OpName, &ti.LogicTime, &ti.EventIndex, &ti.InputData,
		&ti.TxHash, &ti.SendTime, &ti.Status, &ti.BlockNumber)
}

func TestColStr(t *testing.T) {
	expected := "op_index,op_name,logic_time,event_index,input_data,tx_hash,send_time,status,block_number"
	assert.Equal(t, colStr(&TestItem{}), expected)
}

func TestSelectQuery(t *testing.T) {
	expected := "SELECT " +
		"op_index,op_name,logic_time,event_index,input_data,tx_hash,send_time,status,block_number " +
		"FROM testItems " +
		"WHERE op_index=?"
	assert.Equal(t, selectQuery(&TestItem{}, "WHERE op_index=?"), expected)
}

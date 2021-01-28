// Package usql is a collections of helpers for interacting with `database/sql`
// It's not intended to be a ORM

/* Example:

// Each field must can be scanned from the corresponding table
// column (tagged by sqlcol) via RowScanner.Scan()
type Op struct {
	OpIndex    int64     `sqlcol:"op_index"`
	OpName     string    `sqlcol:"op_name"`
	LogicTime  time.Time `sqlcol:"logic_time"`
	EventIndex int64     `sqlcol:"event_index"`
	InputData  string    `sqlcol:"input_data"`

	TxHash      sql.NullString `sqlcol:"tx_hash"`
	SendTime    mysql.NullTime `sqlcol:"send_time"`
	Status      string         `sqlcol:"status"`
	BlockNumber sql.NullInt64  `sqlcol:"block_number"`
}

// implement usql.DBItem interface
func (op *Op) TableName() string {
	return "ops"
}

// implement usql.DBItem interface
func (op *Op) ScanRow(row RowScanner) error {
	return row.Scan(
		&op.OpIndex, &op.OpName, &op.LogicTime, &op.EventIndex, &op.InputData,
		&op.TxHash, &op.SendTime, &op.Status, &op.BlockNumber)
}

// The query starts from the WHERE clause. Returns err usql.NotFound if no such an op.
func SelectOp(stub Querier, query string, args ...interface{}) (*Op, error) {
	op := Op{}
	err := SelectItem(stub, query, &op, args...)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

// The query starts from the WHERE clause.
func SelectOps(stub Querier, query string, args ...interface{}) ([]*Op, error) {
	ops, err := SelectItems(stub, query, &Op{}, args...)
	if err != nil {
		return nil, err
	}

	var res []*Op
	for _, op := range ops {
		res = append(res, op.(*Op))
	}
	return res, nil
}
*/
package usql

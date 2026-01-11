//go:build integration

package integration

import (
	"testing"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/table"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/stretchr/testify/assert"
)

func TestDecimalExp(t *testing.T) {
	dbConn := testx.GlobalEnv().DBConn()

	iResult, err := dbConn.Exec(`
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO decimal_exp
(id, balance_a, balance_b, balance_history)
values
(uuid_generate_v4(), 2849.9873, null, ARRAY[1000.50, 2000.75, 3000.25]::decimal[]),
(uuid_generate_v4(), 2849.1000000, 274629.55, ARRAY[1000.1232, 2000.75, 3000.25]::decimal[])
;
`)
	assert.NoError(t, err)
	rowsAffected, err := iResult.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), rowsAffected)

	stmt := table.DecimalExp.SELECT(table.DecimalExp.AllColumns)

	var result []entity.DecimalExp
	err = stmt.Query(dbConn, &result)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	assert.Equal(t, "2849.9873", result[0].BalanceA.String())
	assert.Nil(t, result[0].BalanceB)
	assert.Equal(t, 3, len(result[0].BalanceHistory))
	assert.Equal(t, "1000.5", result[0].BalanceHistory[0].String())
	assert.Equal(t, "1000.50", result[0].BalanceHistory[0].StringFixed(2))
	assert.Equal(t, int32(-2), result[0].BalanceHistory[0].Exponent())
	assert.Equal(t, "2000.75", result[0].BalanceHistory[1].String())
	assert.Equal(t, "3000.25", result[0].BalanceHistory[2].String())

	assert.Equal(t, "2849.1", result[1].BalanceA.String())
	assert.Equal(t, "2849.1000000", result[1].BalanceA.StringFixed(7))
	assert.Equal(t, int32(-7), result[1].BalanceA.Exponent())
	assert.NotNil(t, result[1].BalanceB)
	assert.Equal(t, "274629.55", result[1].BalanceB.String())
	assert.Equal(t, 3, len(result[1].BalanceHistory))
	assert.Equal(t, "1000.1232", result[1].BalanceHistory[0].String())
	assert.Equal(t, "2000.75", result[1].BalanceHistory[1].String())
	assert.Equal(t, "3000.25", result[1].BalanceHistory[2].String())
}

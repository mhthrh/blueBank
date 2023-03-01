package Test

import (
	_ "github.com/lib/pq"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CreateBilan(t *testing.T) {
	db := Db.NewDb(cnn)
	err := db.NewBilan(ctx, bilan)
	require.NoError(t, err)

	bilan.IsCredit = false
	err = db.NewBilan(ctx, bilan)
	assert.EqualErrorf(t, err, "type mismatch", "an to an")
}

func Test_BalanceBilan(t *testing.T) {
	db := Db.NewDb(cnn)
	amount, err := db.BalanceBilan(ctx, bilan)

	require.NoError(t, err)
	require.Equal(t, amount, bilan.Amount)
}

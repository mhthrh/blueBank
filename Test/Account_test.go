package Test

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Utils/RandomUtil"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createAccount(t *testing.T) {
	_ = createUser()
	db := Db.NewDb(cnn)
	err := db.CreateAccount(ctx, account)
	require.NoError(t, err)

	//account.Balance = 10
	//err1 := db.CreateAccount(ctx, account)
	//assert.Error(t, err1, "just zero amount/lock amount acceted", "an to an")
}

func Test_CreateAccount(t *testing.T) {
	_ = createUser()
	type args struct {
		ctx     context.Context
		account Entity.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				ctx: context.Background(),
				account: Entity.Account{
					CustomerUserName: username,
					AccountNumber:    RandomUtil.RandomInt(10000, 100000),
					Balance:          0,
					LockAmount:       0,
					CreateAt:         time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "test-2",
			args: args{
				ctx: context.Background(),
				account: Entity.Account{
					CustomerUserName: username,
					AccountNumber:    RandomUtil.RandomInt(10000, 100000),
					Balance:          1,
					LockAmount:       0,
					CreateAt:         time.Time{},
				},
			},
			wantErr: true,
		},
		{
			name: "test-3",
			args: args{
				ctx: context.Background(),
				account: Entity.Account{
					CustomerUserName: username,
					AccountNumber:    accountNumber,
					Balance:          0,
					LockAmount:       100,
					CreateAt:         time.Time{},
				},
			},
			wantErr: true,
		},
	}
	dataBase := Db.NewDb(cnn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dataBase.CreateAccount(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_BalanceAccount(t *testing.T) {
	dataBase := Db.NewDb(cnn)
	username := RandomUtil.RandomString(10)
	u := Entity.Customer{
		FullName: RandomUtil.RandomString(10),
		UserName: username,
		PassWord: password,
		Email:    fmt.Sprintf("%s@gmail.com", RandomUtil.RandomString(10)),
	}

	a := RandomUtil.RandomInt(10000, 100000)
	_ = dataBase.Create(ctx, u)
	type args struct {
		ctx     context.Context
		account Entity.Account
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				ctx: context.Background(),
				account: Entity.Account{
					CustomerUserName: username,
					AccountNumber:    a,
					Balance:          0,
					LockAmount:       0,
					CreateAt:         time.Time{},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "test-2",
			args: args{
				ctx: context.Background(),
				account: Entity.Account{
					CustomerUserName: username,
					AccountNumber:    a,
					Balance:          0,
					LockAmount:       0,
					CreateAt:         time.Time{},
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	_ = dataBase.CreateAccount(context.Background(), tests[0].args.account)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dataBase.BalanceAccount(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BalanceAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package dg_yop_test_user

import (
	"context"
	"testing"

	"hotbox-adm-backend/cli"
)

func init() {
	cli.InitEnv()
	// cli.InitGormDB()
	cli.InitHDGormDB()
	cli.InitHDTaskDB()
}

func TestDgYopTestUserIncomeCornJob_StartSumIncome(t *testing.T) {
	p := &DgYopTestUserIncomeCornJob{}
	p.StartSumIncome(context.Background())
}

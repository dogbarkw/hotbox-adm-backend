package daily_gmv

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

func TestDailyGmvCornJob_Run(t *testing.T) {
	p := &DailyGmvCornJob{}
	p.Run()
}

func TestStartStatGmv_Run(t *testing.T) {
	start := "2025-03-17 00:00:00"
	end := "2025-03-18 23:59:59"
	StartStatGmv(context.Background(), start, end)
}

package artist_nft_activity_corn

import (
	"testing"

	"hotbox-adm-backend/cli"

	"github.com/joho/godotenv"
)

func init() {
	paths := []string{
		"./.env",
		"../.env",
		"../../.env",
		"../../../.env",
	}
	var e error
	for _, v := range paths {
		err := godotenv.Load(v)
		e = err
		if err == nil {
			break
		}
	}
	if e != nil {
		panic(e)
	}

	cli.InitEnv()
	cli.InitSpecialUserIds()
	cli.InitGormDB()
	cli.InitHDGormDB()
	cli.InitHDRedis()
	cli.InitHDADBGormDB()
}

func Test_startArtistNftActivityCornJob(t *testing.T) {
	startArtistNftActivityCornJob()
}

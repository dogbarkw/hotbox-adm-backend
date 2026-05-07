package httpReq

import (
	"os"
	"testing"
)

func TestFeiShuRootBot(t *testing.T) {
	os.Setenv("ENV", "qa")
	type args struct {
		format string
		msg    []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				format: "%d 回收 %s (mobile: %s) 如下藏品: 《%s》%d 份",
				msg:    []any{111, "卢本伟", "123456", "老虎", 2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FeiShuRootBot(tt.args.format, tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("FeiShuRootBot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

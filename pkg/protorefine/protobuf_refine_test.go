package protorefine

import (
	"testing"
)

func Test_protobufImpl_Refine(t *testing.T) {
	type fields struct {
		srcDir string
	}
	type args struct {
		pbImportPath string
		protoDir     string
		outputDir    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_protobufImpl_Refine",
			fields: fields{
				srcDir: ".",
			},
			args: args{
				pbImportPath: "payment/autogen/pb",
				protoDir:     "D:\\codes\\hdmall\\common\\proto",
				outputDir:    "./autogen",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New("D:\\codes\\hdmall\\backend\\payment").Refine(tt.args.pbImportPath, tt.args.protoDir, tt.args.outputDir, "autogen"); (err != nil) != tt.wantErr {
				t.Errorf("Refine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

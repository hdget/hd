package sourcecode

import "testing"

func Test_sourceCodeHandlerImpl_Handle(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test_sourceCodeHandlerImpl_Handle",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcDir := "D://Codes//template//backend//base"

			if err := New(srcDir, WithSkipDirs("autogen")).Handle(); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

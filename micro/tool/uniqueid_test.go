package tool

import "testing"

func TestGenId(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenId(); tt.want {
				t.Errorf("GenId() = %v", got)
			}
		})
	}
}
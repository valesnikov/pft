package main

import "testing"

func Test_getBarBySize(t *testing.T) {
	type args struct {
		size      int
		progress  float64
		roundPrec int
	}
	tests := []struct {
		name string
		args args
	}{
		{"0", args{0, 32, 0}},
		{"1", args{1, 73, 1}},
		{"2", args{2, 25, 2}},
		{"3", args{3, 45, 0}},
		{"4", args{4, 06, 1}},
		{"5", args{5, 99, 2}},
		{"6", args{6, 100, 0}},
		{"7", args{7, 2, 1}},
		{"8", args{8, 84, 2}},
		{"9", args{9, 15, 0}},
		{"10", args{10, 38, 1}},
		{"11", args{11, 73, 2}},
		{"12", args{12, 59, 0}},
		{"13", args{13, 86, 1}},
		{"14", args{14, 23, 2}},
		{"15", args{15, 84, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBarBySize(tt.args.size, tt.args.progress, tt.args.roundPrec); len(got) > tt.args.size {
				t.Errorf("getBarBySize() = %v, more then %v", got, tt.args.size)
			}
		})
	}
}

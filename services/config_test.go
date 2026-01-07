package services

import (
	"fmt"
	"testing"
)

func TestNormalizeDashscopeBaseURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		in      string
		want    string
		changed bool
	}{
		{
			name:    "empty",
			in:      "",
			want:    "",
			changed: false,
		},
		{
			name:    "trim_and_remove_trailing_slash",
			in:      " https://dashscope.aliyuncs.com/compatible-mode/v1/ ",
			want:    "https://dashscope.aliyuncs.com/compatible-mode/v1",
			changed: true,
		},
		{
			name:    "fix_compatible_moe_and_dv1",
			in:      "https://dashscope.aliyuncs.com/compatible-moe/dv1",
			want:    "https://dashscope.aliyuncs.com/compatible-mode/v1",
			changed: true,
		},
		{
			name:    "append_v1_when_only_compatible_mode",
			in:      "https://dashscope.aliyuncs.com/compatible-mode",
			want:    "https://dashscope.aliyuncs.com/compatible-mode/v1",
			changed: true,
		},
		{
			name:    "already_ok",
			in:      "https://dashscope.aliyuncs.com/compatible-mode/v1",
			want:    "https://dashscope.aliyuncs.com/compatible-mode/v1",
			changed: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, changed := normalizeDashscopeBaseURL(tc.in)
			if got != tc.want || changed != tc.changed {
				t.Fatalf("normalizeDashscopeBaseURL(%q) = (%q, %v), want (%q, %v)", tc.in, got, changed, tc.want, tc.changed)
			}
		})
	}
}

func TestGetAppDataDir(t *testing.T) {
	fmt.Println(GetAppDataDir())
}

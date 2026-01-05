package services

import "testing"

func TestExtractSectionImpl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		text        string
		startMarker string
		endMarker   string
		want        string
	}{
		{
			name:        "basic_between_markers",
			text:        "摘要AAA基本面分析BBB技术面分析CCC",
			startMarker: "摘要",
			endMarker:   "基本面分析",
			want:        "AAA",
		},
		{
			name:        "start_missing_returns_empty",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "不存在",
			endMarker:   "基本面分析",
			want:        "",
		},
		{
			name:        "end_missing_returns_to_end",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "摘要",
			endMarker:   "不存在",
			want:        "AAA基本面分析BBB",
		},
		{
			name:        "start_empty_from_zero",
			text:        "AAA基本面分析BBB",
			startMarker: "",
			endMarker:   "基本面分析",
			want:        "AAA",
		},
		{
			name:        "end_empty_to_end",
			text:        "摘要AAA基本面分析BBB",
			startMarker: "摘要",
			endMarker:   "",
			want:        "AAA基本面分析BBB",
		},
		{
			name:        "end_before_start_returns_empty",
			text:        "ENDxxxSTARTyyy",
			startMarker: "START",
			endMarker:   "END",
			want:        "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := extractSectionImpl(tc.text, tc.startMarker, tc.endMarker)
			if got != tc.want {
				t.Fatalf("extractSectionImpl() = %q, want %q", got, tc.want)
			}
		})
	}
}



package asciiArt

import (
	"testing"
)

func TestLoadBannerMap(t *testing.T) {
	tests := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "nonexistentfile.txt",
			Expect: "invalid bannerfile name",
		}, {
			Input:  "standard",
			Expect: "banners/standard.txt",
		},
	}

	for _, test := range tests {
		if got := BannerFile(test.Input); got != test.Expect {
			t.Errorf("Expected %s, but got %s\n", test.Expect, got)
		}
	}
}

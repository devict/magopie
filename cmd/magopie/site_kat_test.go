package main

import (
	"os"
	"reflect"
	"testing"

	mp "github.com/devict/magopie"
)

func TestKATParse(t *testing.T) {
	data, err := os.Open("testdata/kat_ubuntu.xml")
	if err != nil {
		t.Fatal("Error opening fixture", err)
	}
	defer data.Close()

	actual, err := katParse(data)
	if err != nil {
		t.Fatal("katParse err should be nil, was:", err)
	}

	expected := []mp.Torrent{
		{
			ID:        "0C5B427C2F833B09EA0E3DC7C624F6C187125267",
			Title:     "Ubuntu Linux Go from Beginner to Power User!",
			MagnetURI: "magnet:?xt=urn:btih:0C5B427C2F833B09EA0E3DC7C624F6C187125267&dn=ubuntu+linux+go+from+beginner+to+power+user&tr=udp%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce",
			SiteID:    "kat",
			Size:      1137244499,
			Seeders:   85,
			Leechers:  105,
		},
		{
			ID:        "21236117B7A773639BD5C7C771E66A045BD51A8A",
			Title:     "Learning Ubuntu Linux Server",
			MagnetURI: "magnet:?xt=urn:btih:21236117B7A773639BD5C7C771E66A045BD51A8A&dn=learning+ubuntu+linux+server&tr=udp%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce",
			SiteID:    "kat",
			Size:      581653341,
			Seeders:   48,
			Leechers:  12,
		},
		{
			ID:        "13DBA979D53F20E6A73D4EE939952D1C367B64C7",
			Title:     "CodeWeavers Crossover 15.0.1 with crack for ubuntu fedora linux",
			MagnetURI: "magnet:?xt=urn:btih:13DBA979D53F20E6A73D4EE939952D1C367B64C7&dn=codeweavers+crossover+15+0+1+with+crack+for+ubuntu+fedora+linux&tr=udp%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce",
			SiteID:    "kat",
			Size:      240330669,
			Seeders:   1,
			Leechers:  1,
		},
		{
			ID:        "EB7B040141407F150E32FF366CD624403387B5C1",
			Title:     "Ubuntu 16.04 (Xenial Xerus) Alpha Desktop AMD64 (64-bit PC)",
			MagnetURI: "magnet:?xt=urn:btih:EB7B040141407F150E32FF366CD624403387B5C1&dn=ubuntu+16+04+xenial+xerus+alpha+desktop+amd64+64+bit+pc&tr=udp%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce",
			SiteID:    "kat",
			Size:      1480048640,
			Seeders:   1,
			Leechers:  3,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tpbParse actual = %v\nexpected %v", actual, expected)
	}
}

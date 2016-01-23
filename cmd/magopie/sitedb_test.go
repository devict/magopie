package main

import (
	"reflect"
	"testing"

	mp "github.com/gophergala2016/magopie"
)

var (
	siteFoo = site{Site: mp.Site{ID: "foo", Enabled: true}}
	siteBar = site{Site: mp.Site{ID: "bar", Enabled: false}}
	siteBaz = site{Site: mp.Site{ID: "baz", Enabled: true}}
)

func TestGetSite(t *testing.T) {
	a := sitedb{
		sites: []site{siteFoo, siteBar},
	}

	if actual := a.GetSite(siteBar.ID); !reflect.DeepEqual(actual, siteBar) {
		t.Errorf("sitedb.GetSite(%q) = %v, expected %v", siteBar.ID, actual, siteBar)
	}
}

func TestGetAllSites(t *testing.T) {
	a := sitedb{
		sites: []site{siteFoo, siteBar, siteBaz},
	}

	expected := []site{siteFoo, siteBar, siteBaz}
	if actual := a.GetAllSites(); !reflect.DeepEqual(actual, expected) {
		t.Errorf("sitedb.GetAllSites() = %v, expected %v", actual, expected)
	}
}

func TestGetEnabledSites(t *testing.T) {
	a := sitedb{
		sites: []site{siteFoo, siteBar, siteBaz},
	}

	expected := []site{siteFoo, siteBaz}
	if actual := a.GetEnabledSites(); !reflect.DeepEqual(actual, expected) {
		t.Errorf("sitedb.GetAllSites() = %v, expected %v", actual, siteBar)
	}
}

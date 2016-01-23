package main

import (
	"reflect"
	"testing"

	"github.com/gophergala2016/magopie/entities"
)

func TestGetSite(t *testing.T) {
	foo := entities.Site{ID: "foo"}
	bar := entities.Site{ID: "bar"}

	a := app{
		sites: []entities.Site{foo, bar},
	}

	if actual := a.GetSite("bar"); actual != bar {
		t.Errorf("app.GetSite(bar) = %v, expected %v", actual, bar)
	}
}

func TestGetAllSites(t *testing.T) {
	foo := entities.Site{ID: "foo", Enabled: true}
	bar := entities.Site{ID: "bar", Enabled: false}
	baz := entities.Site{ID: "baz", Enabled: true}

	a := app{
		sites: []entities.Site{foo, bar, baz},
	}

	expected := []entities.Site{foo, bar, baz}
	if actual := a.GetAllSites(); !reflect.DeepEqual(actual, expected) {
		t.Errorf("app.GetAllSites() = %v, expected %v", actual, bar)
	}
}

func TestGetEnabledSites(t *testing.T) {
	foo := entities.Site{ID: "foo", Enabled: true}
	bar := entities.Site{ID: "bar", Enabled: false}
	baz := entities.Site{ID: "baz", Enabled: true}

	a := app{
		sites: []entities.Site{foo, bar, baz},
	}

	expected := []entities.Site{foo, baz}
	if actual := a.GetEnabledSites(); !reflect.DeepEqual(actual, expected) {
		t.Errorf("app.GetAllSites() = %v, expected %v", actual, bar)
	}
}

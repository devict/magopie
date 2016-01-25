package main

import (
	"sync"

	"github.com/devict/magopie"
)

type site struct {
	magopie.Site
	search func(term string) ([]magopie.Torrent, error)
}

type sitedb struct {
	mtx   sync.RWMutex
	sites []site
}

func (db *sitedb) GetSite(id string) site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	for _, s := range db.sites {
		if s.ID == id {
			return s
		}
	}

	return site{}
}

func (db *sitedb) GetAllSites() []site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	return db.sites
}

func (db *sitedb) GetEnabledSites() []site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	var sites []site
	for _, s := range db.sites {
		if s.Enabled {
			sites = append(sites, s)
		}
	}
	return sites
}

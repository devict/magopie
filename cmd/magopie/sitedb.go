package main

import "github.com/gophergala2016/magopie/entities"
import "sync"

type sitedb struct {
	mtx   sync.RWMutex
	sites []entities.Site
}

func (db *sitedb) GetSite(id string) entities.Site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	for _, s := range db.sites {
		if s.ID == id {
			return s
		}
	}

	return entities.Site{}
}

func (db *sitedb) GetAllSites() []entities.Site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	return db.sites
}

func (db *sitedb) GetEnabledSites() []entities.Site {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	var sites []entities.Site
	for _, s := range db.sites {
		if s.Enabled {
			sites = append(sites, s)
		}
	}
	return sites
}

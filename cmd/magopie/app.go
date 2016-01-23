package main

import "github.com/gophergala2016/magopie/entities"
import "sync"

type app struct {
	mtx   sync.RWMutex
	sites []entities.Site
}

func (a *app) GetSite(id string) entities.Site {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	for _, s := range a.sites {
		if s.ID == id {
			return s
		}
	}

	return entities.Site{}
}

func (a *app) GetAllSites() []entities.Site {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	return a.sites
}

func (a *app) GetEnabledSites() []entities.Site {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	var sites []entities.Site
	for _, s := range a.sites {
		if s.Enabled {
			sites = append(sites, s)
		}
	}
	return sites
}

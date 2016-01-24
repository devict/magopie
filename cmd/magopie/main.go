package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/spf13/afero"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	a := &server{
		fs: afero.OsFs{},
		sites: sitedb{
			sites: []site{
				kickAssTorrents,
				thePirateBay,
			},
		},
		torcacheURL: "http://torcache.net",
		downloadDir: ".",
	}

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", "localhost", port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router(a)))
}

func router(a *server) http.Handler {
	mux := xmux.New()

	c := xhandler.Chain{}

	c.Use(mwLogger)

	mux.GET("/sites", c.HandlerCF(xhandler.HandlerFuncC(a.handleAllSites)))
	mux.GET("/sites/:id", c.HandlerCF(xhandler.HandlerFuncC(a.handleSingleSite)))
	mux.GET("/torrents", c.HandlerCF(xhandler.HandlerFuncC(a.handleTorrents)))
	mux.POST("/download/:hash", c.HandlerCF(xhandler.HandlerFuncC(a.handleDownload)))

	return xhandler.New(context.Background(), mux)
}

func mwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving", r.Method, r.URL.String(), "to", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

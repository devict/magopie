package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devict/magopie"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/spf13/afero"
)

var (
	addr        = flag.String("addr", ":8080", "HTTP(S) service address")
	downloadDir = flag.String("dir", ".", "Directory where magopie should download .torrent files")
	apiKey      = flag.String("key", "", "Shared API key for clients (required)")
	tlsKey      = flag.String("tlsKey", "", "Path to TLS Key")
	tlsCert     = flag.String("tlsCert", "", "Path to TLS Cert")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "-key flag is required")
		flag.Usage()
		os.Exit(1)
	}

	a := &server{
		key: *apiKey,
		fs:  afero.OsFs{},
		sites: sitedb{
			sites: []site{
				kickAssTorrents,
				thePirateBay,
			},
		},
		torcacheURL: "http://torcache.net",
		downloadDir: *downloadDir,
	}

	if *tlsKey != "" && *tlsCert != "" {
		log.Printf("Listening for HTTPS on %s with key %s and cert %s", *addr, *tlsKey, *tlsCert)
		log.Fatal(http.ListenAndServeTLS(*addr, *tlsCert, *tlsKey, router(a)))
	} else {
		log.Printf("Listening for HTTP on %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, router(a)))
	}
}

func router(a *server) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/sites", a.handleAllSites).Methods("GET")
	r.HandleFunc("/sites/{id}", a.handleSingleSite).Methods("GET")
	r.HandleFunc("/torrents", a.handleTorrents).Methods("GET")
	r.HandleFunc("/download/{hash}", a.handleDownload).Methods("POST")

	chain := alice.New(mwLogger, mwAuthenticationCheck(a.key)).Then(r)

	return chain
}

func mwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving", r.Method, r.URL.String(), "to", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func mwAuthenticationCheck(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !magopie.CheckMAC(r.Header.Get("X-Request-ID"), r.Header.Get("X-HMAC"), key) {
				log.Println("Request failed HMAC")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

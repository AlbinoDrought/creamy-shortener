package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/AlbinoDrought/creamy-shortener/linker"
	"github.com/AlbinoDrought/creamy-shortener/linker/guard/hostguard"
	"github.com/AlbinoDrought/creamy-shortener/linker/mapper/multihashmapper"
	"github.com/AlbinoDrought/creamy-shortener/linker/repo/aferorepo"

	"github.com/spf13/afero"

	"github.com/gorilla/mux"
)

var appURL string
var port string
var localDataPath string
var hosts []string
var hashMode string

var linkRepo *linker.LinkShortener

func envDefault(name string, backup string) string {
	found, exists := os.LookupEnv(name)
	if exists {
		return found
	}
	return backup
}

func init() {
	appURL = envDefault("CREAMY_APP_URL", "http://localhost:3000/")
	hashMode = envDefault("CREAMY_HASH_MODE", "sha2-256")
	port = envDefault("CREAMY_HTTP_PORT", "3000")
	localDataPath = envDefault("CREAMY_DATA_PATH", "./data")
	hosts = strings.Split(envDefault("CREAMY_POPULATED_HOSTS", "localhost"), ",")
}

func viewLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["link"]

	link, err := linkRepo.Expand(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("link not found"))
		return
	}

	http.Redirect(w, r, link.String(), http.StatusPermanentRedirect)
}

func shortenLink(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("bad data"))
		return
	}

	link := r.FormValue("link")
	if link == "" {
		w.WriteHeader(422)
		w.Write([]byte("missing \"link\" value"))
		return
	}

	parsedLink, err := url.Parse(link)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("bad \"link\" value"))
		return
	}

	if !linkRepo.Allowed(parsedLink) {
		w.WriteHeader(422)
		w.Write([]byte("link not allowed"))
		return
	}

	id, err := linkRepo.Shorten(parsedLink)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("unhandled error shortening link"))
		return
	}

	fullURL, _ := url.Parse(appURL)
	fullURL.Path = "/l/" + id

	w.Write([]byte(fullURL.String()))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am a talking link shortener"))
}

func main() {
	_, err := url.Parse(appURL)
	if err != nil {
		log.Fatal("unable to parse app url", appURL, err)
	}

	linkRepo = linker.Make(
		hostguard.Make(hosts),
		multihashmapper.Must(hashMode),
		aferorepo.Must(afero.NewBasePathFs(afero.NewOsFs(), localDataPath)),
	)

	router := makeRouter([]routeDef{
		routeDef{"GET", "/l/{link}", "ViewLink", viewLink},
		routeDef{"POST", "/shorten", "ShortenLink", shortenLink},
		routeDef{"GET", "/", "Welcome", welcome},
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

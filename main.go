package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var appURL string
var port string
var localDataPath string
var hosts []string
var hashMode string

var linkRepo linker

func envDefault(name string, backup string) string {
	found, exists := os.LookupEnv(name)
	if exists {
		return found
	}
	return backup
}

func init() {
	appURL = envDefault("CREAMY_APP_URL", "")
	hashMode = envDefault("CREAMY_HASH_MODE", "sha2-256")
	port = envDefault("CREAMY_HTTP_PORT", "3000")
	localDataPath = envDefault("CREAMY_DATA_PATH", "./data")
	hosts = strings.Split(envDefault("CREAMY_POPULATED_HOSTS", "localhost"), ",")
}

func viewLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortenedPart := vars["link"]

	link, err := linkRepo.Expand(shortenedPart)
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

	if err = linkRepo.Allowed(parsedLink); err != nil {
		w.WriteHeader(422)
		w.Write([]byte("link not allowed"))
		return
	}

	shortened, err := linkRepo.Shorten(parsedLink)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("unhandled error shortening link"))
		return
	}

	w.Write([]byte(shortened))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am a talking link shortener"))
}

func main() {
	linkRepo = makeLinker(localDataPath, appURL, hashMode)

	for _, host := range hosts {
		parsedHost, err := url.Parse("https://" + host)
		if err != nil {
			log.Fatal("could not parse allowed host", host, err)
		}

		if err := linkRepo.Allow(parsedHost.Host); err != nil {
			log.Fatal("failed to allow host", host, err)
		}

		if err := linkRepo.Allowed(parsedHost); err != nil {
			log.Fatal("tried to allow host but it didn't work", host, err)
		}
	}

	router := makeRouter([]routeDef{
		routeDef{"GET", "/l/{link}", "ViewLink", viewLink},
		routeDef{"POST", "/shorten", "ShortenLink", shortenLink},
		routeDef{"GET", "/", "Welcome", welcome},
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

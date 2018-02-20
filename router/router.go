package router

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/sensiblecodeio/tiny-ssl-reverse-proxy/pkg/wsproxy"
)

var u *url.URL
var usp string
var ust http.RoundTripper
var nc *http.Client
var tls bool

type Config struct {
	GitlabSocket string
	Uri          string
	TLS          bool
}

type Route struct {
	Path                string
	Method              string
	Auth                bool
	HandlerStruct       interface{}
	HandlerMethod       string
	HandlerQueryStrings []string
}

type RouterAPI struct {
	router *mux.Router
}

func New(rcfg Config) *mux.Router {
	var err error
	u, err = url.Parse(rcfg.Uri)
	if err != nil {
		log.Fatal(err.Error())
	}
	usp = rcfg.GitlabSocket
	ust = &http.Transport{Dial: unixSocketDial}
	tls = rcfg.TLS
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(proxyGitlab)
	return r
}

func NewRouterAPI(r *mux.Router) RouterAPI {
	nc = &http.Client{Timeout: time.Second * 10, Transport: ust}
	return RouterAPI{router: r}
}

func (r *RouterAPI) AddRoute(route Route) {
	r.router.Path(route.Path).Methods(route.Method).HandlerFunc(func(route Route) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			tStart := time.Now()
			var err error
			if route.Auth {
				err = r.requireAuth(*req)
			}
			if tls {
				ensureSTS(w)
			}
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				errorWriter(w, http.StatusInternalServerError, err)
				return
			} else {
				inputs := []reflect.Value{}
				inputs = append(inputs, reflect.ValueOf(mux.Vars(req)))
				queries := req.URL.Query()
				inputs = append(inputs, reflect.ValueOf(queries))
				out := reflect.ValueOf(route.HandlerStruct).MethodByName(route.HandlerMethod).Call(inputs)
				resp := out[0].Interface()
				if out[len(out)-1].Interface() != nil {
					err = out[len(out)-1].Interface().(error)
				}
				tEnd := time.Now()
				t := tEnd.Sub(tStart)
				if err != nil {
					errorWriter(w, http.StatusInternalServerError, err)
					return
				}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					errorWriter(w, http.StatusInternalServerError, err)
					return
				}
				metadataToHeaders(w, out[1].Interface().(map[string]string))
				w.Header().Set("X-Runtime", t.String())
				if route.Method == "POST" {
					w.WriteHeader(http.StatusCreated)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Write(jsonResp)
			}
		}
	}(route))
}

func ensureSTS(w http.ResponseWriter) {
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
}

func metadataToHeaders(w http.ResponseWriter, metadata map[string]string) {
	for key, header := range metadata {
		w.Header().Set("X-"+key, header)
	}
}

func (r *RouterAPI) requireAuth(req http.Request) error {
	nr, err := http.NewRequest("GET", u.String()+"/api/v4/version", nil)
	if err != nil {
		panic(err)
	}
	nr.Header = req.Header
	resp, err := nc.Do(nr)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		err = errors.New("401 Unauthorized")
	} else {
		err = nil
	}
	return err
}

func errorWriter(w http.ResponseWriter, statusCode int, e error) {
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"message\": \"" + e.Error() + "\"}"))
}

func unixSocketDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", usp)
}

func proxyGitlab(w http.ResponseWriter, r *http.Request) {
	if tls {
		ensureSTS(w)
		r.Header.Set("X-Forwarded-Proto", "https")
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.Transport = &http.Transport{Dial: unixSocketDial}

	(&wsproxy.ReverseProxy{rp}).ServeHTTP(w, r)
}

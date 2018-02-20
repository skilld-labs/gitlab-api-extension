package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"./api"
	"./api/v4"
	"./db"
	"./router"
)

func main() {
	bindAddress := flag.String("bindAddress", "", "The address (incl. port) to bind")
	dbSocketPath := flag.String("dbSocketPath", "", "The database socket path")
	dbName := flag.String("dbName", "", "The database name")
	gitlabSocketPath := flag.String("gitlabSocketPath", "", "The gitlab workhorse socket path")
	uri := flag.String("uri", "", "The gitlab instance uri")
	tlsCertificate := flag.String("tlsCertificate", "", "The TLS certificate file")
	tlsKey := flag.String("tlsKey", "", "The TLS key file")
	flag.Parse()

	d, err := db.New(db.Config{SocketPath: *dbSocketPath, Name: *dbName})
	if err != nil {
		log.Fatal(err.Error())
	}
	dapi := db.NewDbAPI(d)

	a := api.New(api.Config{DbAPI: *dapi})
	aapiV4 := apiv4.NewApiAPI(a)

	enableTls := *tlsCertificate != "" || *tlsKey != ""
	r := router.New(router.Config{GitlabSocket: *gitlabSocketPath, Uri: *uri, TLS: enableTls})
	rapi := router.NewRouterAPI(r)
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetTimelogs",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/projects/{projectID}/issues/{issueIID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetIssueTimelogs",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/projects/{projectID}/merge_requests/{mergeRequestIID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetMergeRequestTimelogs",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/users/{userID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetUserTimelogs",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/projects/{projectID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetProjectTimelogs",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/projects/{projectID}/users/{userID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetUserTimelogsByProject",
		Auth:          true,
	})
	rapi.AddRoute(router.Route{
		Path:          "/api/v4/projects/{projectID}/issues/{issueIID}/users/{userID}/time_logs",
		Method:        "GET",
		HandlerStruct: &aapiV4,
		HandlerMethod: "GetUserTimelogsByProjectAndIssue",
		Auth:          true,
	})

	srv := &http.Server{
		Handler: r,
		Addr:    *bindAddress,
	}

	if enableTls {
		srv.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		log.Fatal(srv.ListenAndServeTLS(*tlsCertificate, *tlsKey))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}

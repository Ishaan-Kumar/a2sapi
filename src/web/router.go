package web

// router.go - request router

import (
	"net/http"
	"strings"
	"time"

	"github.com/syncore/a2sapi/src/config"

	"github.com/syncore/a2sapi/src/logger"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	for _, ar := range apiRoutes {
		handler := http.TimeoutHandler(compressGzip(ar.handlerFunc, config.Config.WebConfig.CompressResponses),
			time.Duration(config.Config.WebConfig.APIWebTimeout)*time.Second,
			`{"error": {"code": 503,"message": "Request timeout."}}`)
		handler = logger.LogWebRequest(handler, ar.name)

		r.Methods(ar.method).
			MatcherFunc(pathQStrToLowerMatcherFunc(r, ar.path, ar.queryStrings,
				getRequiredQryStringCount(ar.queryStrings))).
			Name(ar.name).
			Handler(handler)
	}
	return r
}

// Provide case-insensitive matching for URL paths and query strings
func pathQStrToLowerMatcherFunc(router *mux.Router,
	routepath string, querystrings []querystring,
	requiredQsCount int) func(req *http.Request,
	rt *mux.RouteMatch) bool {
	return func(req *http.Request, rt *mux.RouteMatch) bool {
		pathok, qstrok := false, false
		// case-insensitive paths
		if strings.HasPrefix(strings.ToLower(req.URL.Path), strings.ToLower(routepath)) {
			logger.WriteDebug("PATH: %s matches route path: %s", req.URL.Path, routepath)
			pathok = true
		}
		//case-insensitive query strings
		// not all API routes will make use of query strings
		if len(querystrings) == 0 {
			qstrok = true
		} else {
			qry := req.URL.Query()
			truecount := 0
			for key := range qry {
				logger.WriteDebug("URL query string key is: %s", key)
				for _, qs := range querystrings {
					if strings.EqualFold(key, qs.name) && qs.required {
						logger.WriteDebug("KEY: %s matches query string: %s", key, qs.name)
						truecount++
						break
					}
				}
			}
			if truecount == requiredQsCount {
				qstrok = true
			}
		}
		return pathok && qstrok
	}
}

func getRequiredQryStringCount(querystrings []querystring) int {
	reqcount := 0
	for _, q := range querystrings {
		if q.required {
			reqcount++
		}
	}
	return reqcount
}

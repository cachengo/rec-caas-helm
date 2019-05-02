// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"net/http"
	"pkg/config"

	"github.com/gorilla/mux"
	"github.com/ncw/swift"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*swift.Connection, config.EnvConfig) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"ObjectGet",
		"GET",
		// tonyaw: make it can match slash, such as "aaa/bbb"
		"/{id:.*}",
		ObjectGet,
	},
	Route{
		"ObjectPut",
		"POST",
		// tonyaw: make it can match slash, such as "aaa/bbb"
		"/{id:.*}",
		ObjectPut,
	},
	Route{
		"ObjectPut",
		"PUT",
		// tonyaw: make it can match slash, such as "aaa/bbb"
		"/{id:.*}",
		ObjectPut,
	},
	Route{
		"ObjectDelete",
		"DELETE",
		"/{id:.*}",
		ObjectDelete,
	},
}

func NewRouter(c *swift.Connection, config config.EnvConfig) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(c, config)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

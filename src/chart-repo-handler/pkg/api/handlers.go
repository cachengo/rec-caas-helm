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
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"pkg/config"
	"pkg/repo"

	"github.com/gorilla/mux"
	"github.com/ncw/swift"
)

func Index(c *swift.Connection, config config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>")
		fmt.Fprintln(w, "Hello World!<br />")
		objects, err := c.ObjectNames(config.Container, nil)
		if err != nil {
			log.Println("Index error:", err)
			fmt.Fprintln(w, "Index error:", err)
		} else {
			for _, v := range objects {
				fmt.Fprintf(w, "<a href=\"%s\">%s</a><br />", v, v)
			}
		}
		fmt.Fprintln(w, "</body></html>")
	}
}

func ObjectGet(c *swift.Connection, config config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objectName := vars["id"]

		object, _, err := c.ObjectOpen(config.Container, objectName, false, nil)

		if err != nil {
			switch err.Error() {
			case "Object Not Found":
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("ObjectGet error:", err, ", requested object:", objectName)
			fmt.Fprintln(w, "<html><body>")
			fmt.Fprintln(w, "ObjectGet error:", err)
			fmt.Fprintln(w, "</body></html>")
			return
		}
		io.Copy(w, object)
		object.Close()
	}
}

func ObjectPut(c *swift.Connection, config config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objectName := vars["id"]

		_, err := c.ObjectPut(config.Container, objectName, r.Body, false, "", "", nil)
		if err != nil {
			log.Println("ObjectPut error:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "<html><body>")
			fmt.Fprintln(w, "ObjectPut error:", err)
			fmt.Fprintln(w, "</body></html>")
			return
		}
		if strings.HasPrefix(objectName, config.IndexPath) && strings.HasSuffix(objectName, "tgz") {
			log.Println("Regenerating index for:", config.IndexPath)
			if err := repo.Index(c, config.Container, config.RepoUrl+":"+config.ListenOnPort, config.IndexPath, "/index.yaml"); err != nil {
				log.Println("ObjectPut error:", err.Error())
				fmt.Fprintln(w, "<html><body>")
				fmt.Fprintln(w, "Index file generation error:", err)
				fmt.Fprintln(w, "</body></html>")
			}
		}
	}
}

func ObjectDelete(c *swift.Connection, config config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objectName := vars["id"]

		err := c.ObjectDelete(config.Container, objectName)
		if err != nil {
			log.Println("ObjectDelete error:", err.Error())
			fmt.Fprintln(w, "<html><body>")
			fmt.Fprintln(w, "ObjectDelete error:", err)
			fmt.Fprintln(w, "</body></html>")
			return
		}
		if strings.HasPrefix(objectName, config.IndexPath) && strings.HasSuffix(objectName, "tgz") {
			log.Println("Regenerating index for:", config.IndexPath)
			if err := repo.Index(c, config.Container, config.RepoUrl+":"+config.ListenOnPort, config.IndexPath, "/index.yaml"); err != nil {
				log.Println("ObjectPut error:", err.Error())
				fmt.Fprintln(w, "<html><body>")
				fmt.Fprintln(w, "Index file generation error:", err)
				fmt.Fprintln(w, "</body></html>")
			}
		}
	}
}

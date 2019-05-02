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

// tonyaw: refer to "kubernetes/helm/pkg/repo/index.go"
package repo

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"

	// tonyaw: added import
	"log"
	"pkg/chartutil"

	"github.com/ncw/swift"
	"k8s.io/helm/pkg/provenance"
	helm_repo "k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/urlutil"
)

//var IndexPath = "index.yaml"

// LoadIndexFile takes a file at the given path and returns an IndexFile object
func LoadIndexFile(c *swift.Connection, container string, path string) (*IndexFile, error) {
	// TBD.tonyaw: didin't tested it
	b, err := c.ObjectGetBytes(container, path)
	if err != nil {
		return nil, err
	}
	return loadIndex(b)
}

// tonyaw: return "IndexFile" defined in this file.
func NewIndexFile() *IndexFile {
	var indexFile *helm_repo.IndexFile

	indexFile = &helm_repo.IndexFile{
		APIVersion: helm_repo.APIVersionV1,
		Generated:  time.Now(),
		Entries:    map[string]helm_repo.ChartVersions{},
		PublicKeys: []string{},
	}

	return &IndexFile{
		IndexFile: indexFile,
	}
}

// IndexFile represents the index file in a chart repository
// tonyaw: here only add "WriteObject" function.
type IndexFile struct {
	*helm_repo.IndexFile
}

// WriteObject writes an index file to the given destination path.
func (i IndexFile) WriteObject(c *swift.Connection, container, path string) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	// tonyaw: if "index" exists, this step will overwrite it.
	return c.ObjectPutBytes(container, path, b, "text/plain")
}

// tonyaw: Merge "IndexFile" defined in this file.
func (i *IndexFile) Merge(f *IndexFile) {
	for _, cvs := range f.Entries {
		for _, cv := range cvs {
			if !i.Has(cv.Name, cv.Version) {
				e := i.Entries[cv.Name]
				i.Entries[cv.Name] = append(e, cv)
			}
		}
	}
}

// IndexDirectory reads a (flat) directory and generates an index.
//
// It indexes only charts that have been packaged (*.tgz).
//
// The index returned will be in an unsorted state
func IndexDirectory(c *swift.Connection, container, baseURL, rootPath, indexName string) (*IndexFile, error) {
	// tonyaw: Usw Swift command to replace FileSystem calls
	objects, err := c.ObjectNames(container, nil)

	if err != nil {
		return nil, err
	}
	index := NewIndexFile()
	for _, object := range objects {
		if !strings.HasPrefix(object, rootPath) {
			continue
		}
		if object == filepath.Join(rootPath, indexName) {
			//log.Println("Ignore", filepath.Join(rootPath, indexName))
			continue
		}
		if !strings.HasSuffix(object, ".tgz") {
			log.Println("Ignore", object)
			continue
		}

		//fname := filepath.Base(object)
		chartData, err := chartutil.Load(c, container, object)
		if err != nil {
			log.Println("Error with object:", object, ",", err)
			continue
		}
		//huszty: open the object here for sha256 calculation
		objectReader, _, err := c.ObjectOpen(container, object, false, nil)
		//hash, err := provenance.DigestFile(object)
		hash, err := provenance.Digest(objectReader)
		if err != nil {
			return index, err
		}
		log.Println("Added:", object)
		//huszty: index.Add would split the path from the object if url is provided
		//so join the url with the prefixed object path to fake subdirectories correctly
		objectPath, _ := urlutil.URLJoin(baseURL, object)
		index.Add(chartData.Metadata, objectPath, "", hash)
		objectReader.Close()
	}
	return index, nil
}

// loadIndex loads an index file and does minimal validity checking.
//
// This will fail if API Version is not set (ErrNoAPIVersion) or if the unmarshal fails.
func loadIndex(data []byte) (*IndexFile, error) {
	i := &IndexFile{}
	if err := yaml.Unmarshal(data, i); err != nil {
		return i, err
	}

	return i, nil
}

// tonyaw: refer to "index" from "helm/cmd/helm/repo_index.go".
func Index(c *swift.Connection, container, url, rootPath, indexName string) error {
	i, err := IndexDirectory(c, container, url, rootPath, indexName)
	if err != nil {
		return err
	}
	i.SortEntries()
	return i.WriteObject(c, container, filepath.Join(rootPath, indexName))
}

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

// tonyaw: refer to "kubernetes/helm/pkg/chartutil/load.go"
package chartutil

import (
	"k8s.io/helm/pkg/proto/hapi/chart"

	// tonyaw:
	helm_chartutil "k8s.io/helm/pkg/chartutil"
	"github.com/ncw/swift"
)

// Load takes a string name, tries to resolve it to a file or directory, and then loads it.
//
// This is the preferred way to load a chart. It will discover the chart encoding
// and hand off to the appropriate chart reader.
//
// If a .helmignore file is present, the directory loader will skip loading any files
// matching it. But .helmignore is not evaluated when reading out of an archive.
func Load(c *swift.Connection, container, name string) (*chart.Chart, error) {
	// tonyaw.TBD: ignore dir case for now.
	// fi, err := os.Stat(name)
	// if err != nil {
	// 	return nil, err
	// }
	// if fi.IsDir() {
	// 	return LoadDir(name)
	// }
	return LoadFile(c, container, name)
}


// LoadFile loads from an archive file.
func LoadFile(c *swift.Connection, container, name string) (*chart.Chart, error) {
	// tonyaw.TBD: ignore dir case for now.
	// if fi, err := os.Stat(name); err != nil {
	// 	return nil, err
	// } else if fi.IsDir() {
	// 	return nil, errors.New("cannot load a directory")
	// }

	object, _, err := c.ObjectOpen(container, name, false, nil)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return helm_chartutil.LoadArchive(object)
}


// LoadDir loads from a directory.
//
// This loads charts only from directories.
// tonyaw.TBD: do we still need "LoadDir" for swift?
func LoadDir(dir string) (*chart.Chart, error) {
	return nil, nil
}

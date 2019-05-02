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

package config

type EnvConfig struct {
	AuthUser     string `required:"true"`
	AuthKey      string `required:"true"`
	AuthUrl      string `required:"true"`
	Container    string `required:"true"`
	ListenOnIP   string `required:"true"`
	ListenOnPort string `required:"true"`
	RepoUrl      string `required:"true"`
	IndexPath    string
	TlsCertPath  string
	TlsKeyPath   string
	TlsCaPath   string
}

func (c *EnvConfig) ToString() string {
	result := "AuthUser=" + c.AuthUser
	result += " AuthKey=***"
	result += " AuthUrl=" + c.AuthUrl
	result += " Container=" + c.Container
	result += " ListenOnIP=" + c.ListenOnIP
	result += " ListenOnPort=" + c.ListenOnPort
	result += " RepoUrl=" + c.RepoUrl
	result += " IndexPath=" + c.IndexPath
	result += " TlsCertPath=" + c.TlsCertPath
	result += " TlsKeyPath=" + c.TlsKeyPath
	result += " TlsCaPath" + c.TlsCaPath

	return result
}


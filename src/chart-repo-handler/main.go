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

package main

import (
  "crypto/tls"
  "log"
  "net/http"
  "pkg/api"
  "pkg/config"
  "pkg/repo"
  "github.com/kelseyhightower/envconfig"
  "github.com/ncw/swift"
  "os"
  "os/signal"
  "syscall"
  "strconv"
  "time"
  "errors"
  "crypto/x509"
  "io/ioutil"
)

var (
  chartRepo *http.Server
)

const (
  regenRetryCounter = 10
)

func main() {
  var envConfig config.EnvConfig
  err := envconfig.Process("chartrepohandler", &envConfig)
  if err != nil {
    log.Fatal(err.Error())
    }

  swiftCon := connectSwift(envConfig)
  log.Println("Chart repo handler v0.9 started up.")
  log.Printf("Config: %s\n", envConfig.ToString())
  log.Println("Regenerating index for:", envConfig.IndexPath)
  if err := regenerateIndexYaml(envConfig, &swiftCon); err != nil {
    log.Println(err.Error())
    os.Exit(-1)
    }
  startHttpServer(envConfig, &swiftCon)
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
  for {
    select {
      case <-signalChan:
        log.Println("Shutdown signal received, exiting...")
        if (chartRepo != nil) {
          chartRepo.Close()
          }
        os.Exit(0)
    }
  }
}

func connectSwift(envConfig config.EnvConfig) swift.Connection {
  swiftCon := swift.Connection{
    UserName: 	envConfig.AuthUser,
    ApiKey:   	envConfig.AuthKey,
    AuthUrl:  	envConfig.AuthUrl,
    }
  if envConfig.TlsCaPath != "" {
    log.Printf("INFO: TlsCaPath is presented: Trying to enforce TLS Authentication on swift backend using the server certs")
    file, err := ioutil.ReadFile(envConfig.TlsCaPath)
    if err != nil {
      log.Fatal(err)
      log.Printf("Wrong or missing value in paramteter: TlsCaPath")
      os.Exit(-1)
      }
    certPool := x509.NewCertPool()
    ok := certPool.AppendCertsFromPEM([]byte(file))
    if !ok {
      log.Fatal("Corrupt CACert file")
      os.Exit(-1)
      }
    cert, err := tls.LoadX509KeyPair(envConfig.TlsCertPath, envConfig.TlsKeyPath)
    if err != nil {
      log.Fatal(err)
      log.Printf("Wrong or missing value in paramteter: TlsCertPath or TlsKeyPath")
      os.Exit(-1)
      }
    swiftCon.Transport = &http.Transport{
      TLSClientConfig: &tls.Config{
        RootCAs: certPool,
        Certificates: []tls.Certificate{cert}, 
        },
      }
    }
  return swiftCon
}

func regenerateIndexYaml(envConfig config.EnvConfig, swiftConn *swift.Connection) error {
  var err error
  for i := 0; i <= regenRetryCounter; i++ {
    err = repo.Index(swiftConn, envConfig.Container, envConfig.RepoUrl+":"+envConfig.ListenOnPort, envConfig.IndexPath, "/index.yaml")
    if err != nil {
      log.Println("INFO: Regenerating index.yaml in Swift try no.: " + strconv.Itoa(i) + " was unsuccessful at: " + swiftConn.AuthUrl +
                  " with error:" + err.Error())
      time.Sleep(10 * time.Second)
      continue
    } else {
      return nil
    }
  }
  return errors.New("ERROR: Swift is not responsive, giving-up regeneration!")
}

func startHttpServer(envConfig config.EnvConfig, swiftConn *swift.Connection) {
  router := api.NewRouter(swiftConn, envConfig)
  var tlsCfg *tls.Config
  if envConfig.TlsCertPath != "" && envConfig.TlsKeyPath != "" {
    log.Println("TLS used")
    tlsCfg = &tls.Config{
      MinVersion:               tls.VersionTLS12,
      CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
      PreferServerCipherSuites: true,
      CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
        tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_RSA_WITH_AES_256_CBC_SHA,
      },
    }
  }
  chartRepo = &http.Server{
    Addr:         envConfig.ListenOnIP + ":" + envConfig.ListenOnPort,
    Handler:      router,
    TLSConfig:    tlsCfg,
    TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
  }
  go func() {
    if chartRepo.TLSConfig != nil {
      log.Fatal(chartRepo.ListenAndServeTLS(envConfig.TlsCertPath, envConfig.TlsKeyPath))
    } else {
      log.Fatal(chartRepo.ListenAndServe())
    }
  }()
}


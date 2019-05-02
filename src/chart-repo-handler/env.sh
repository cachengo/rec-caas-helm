#!/bin/sh
# Copyright 2019 Nokia
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export CHARTREPOHANDLER_AUTHKEY=58a2c1331a572c1b69c4
export CHARTREPOHANDLER_AUTHUSER=admin:admin
export CHARTREPOHANDLER_AUTHURL=http://172.24.16.101:8081/auth/v1.0/
export CHARTREPOHANDLER_CONTAINER=packages
export CHARTREPOHANDLER_LISTENONPORT=8088
export CHARTREPOHANDLER_LISTENONIP=0.0.0.0
export CHARTREPOHANDLER_REPOURL=chart-repo.nokia.net
export CHARTREPOHANDLER_INDEXPATH=charts
export CHARTREPOHANDLER_TLSCERTPATH=/etc/etcd/ssl/etcd1.pem
export CHARTREPOHANDLER_TLSKEYPATH=/etc/etcd/ssl/etcd1-key.pem
export CHARTREPOHANDLER_TLSCAPATH=/etc/chart-repo/cacert.pem

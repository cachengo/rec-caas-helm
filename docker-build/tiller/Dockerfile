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

FROM alpine:3.9
MAINTAINER Krisztian Lengyel <krisztian.lengyel@nokia.com>

COPY mainstart.sh /usr/bin/mainstart.sh
COPY tiller /usr/bin/tiller

RUN adduser -u 144 -D -H -s /sbin/nologin kube \
&&  mkdir -p /home/kube \
&&  chown -R kube:kube /home/kube \
&&  chmod a+x /usr/bin/tiller /usr/bin/mainstart.sh

ENTRYPOINT ["/usr/bin/mainstart.sh"]

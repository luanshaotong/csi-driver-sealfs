# Copyright 2020 The Kubernetes Authors.
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

FROM debian:bullseye-20221205

ARG ARCH
ARG binary=./bin/${ARCH}/sealfsplugin

RUN sed -i.bak 's|deb.debian.org|mirrors.tuna.tsinghua.edu.cn|g' /etc/apt/sources.list && \
    apt update && apt upgrade -y && apt-mark unhold libcap2 && \
    apt install -y libfuse3-3 libibverbs1 && \
    apt clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY sealfs-client /usr/bin/sealfs-client
COPY ${binary} /sealfsplugin

ENTRYPOINT ["/sealfsplugin"]

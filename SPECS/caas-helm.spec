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

%define COMPONENT helm
%define RPM_NAME caas-%{COMPONENT}
%define RPM_MAJOR_VERSION 2.14.2
%define RPM_MINOR_VERSION 1
%define IMAGE_TAG %{RPM_MAJOR_VERSION}-%{RPM_MINOR_VERSION}
%define go_version 1.11.5
%define cni_plugins_version 0.7.0
%define binary_build_dir %{_builddir}/%{RPM_NAME}-%{RPM_MAJOR_VERSION}/binary-save
%define docker_build_dir %{_builddir}/%{RPM_NAME}-%{RPM_MAJOR_VERSION}/docker-build
%define docker_save_dir %{_builddir}/%{RPM_NAME}-%{RPM_MAJOR_VERSION}/docker-save
%define built_binaries_dir /binary-save

Name:           %{RPM_NAME}
Version:        %{RPM_MAJOR_VERSION}
Release:        %{RPM_MINOR_VERSION}%{?dist}
Summary:        Containers as a Service %{COMPONENT} component
License:        %{_platform_licence} and MIT license and BSD and Apache License and Lesser General Public License
BuildArch:      x86_64
Vendor:         %{_platform_vendor} and helm/helm unmodified
Source0:        %{name}-%{version}.tar.gz

Requires: docker-ce >= 18.09.2, rsync
BuildRequires: docker-ce-cli >= 18.09.2, rsync, xz

%description
This rpm contains the %{COMPONENT} container for CaaS subsystem.
This container contains the %{COMPONENT} service.

%prep
%autosetup

%build
# Build Helm binaries
docker build \
  --network=host \
  --no-cache \
  --force-rm \
  --build-arg HTTP_PROXY="${http_proxy}" \
  --build-arg HTTPS_PROXY="${https_proxy}" \
  --build-arg NO_PROXY="${no_proxy}" \
  --build-arg http_proxy="${http_proxy}" \
  --build-arg https_proxy="${https_proxy}" \
  --build-arg no_proxy="${no_proxy}" \
  --build-arg HELM_VERSION="%{version}" \
  --build-arg go_version="%{go_version}" \
  --build-arg binaries="%{built_binaries_dir}" \
  --tag helm-builder:%{IMAGE_TAG} \
  %{docker_build_dir}/helm-builder

mkdir -p %{binary_build_dir}
builder_container=$(docker run -id --rm --network=none --entrypoint=/bin/sh helm-builder:%{IMAGE_TAG})
docker cp ${builder_container}:%{built_binaries_dir}/helm   %{binary_build_dir}/
docker cp ${builder_container}:%{built_binaries_dir}/tiller %{binary_build_dir}/
docker rm -f ${builder_container}
docker rmi helm-builder:%{IMAGE_TAG}

# Build tiller container image
rsync -av %{binary_build_dir}/* %{docker_build_dir}/tiller/
docker build \
  --network=host \
  --no-cache \
  --force-rm \
  --build-arg HTTP_PROXY="${http_proxy}" \
  --build-arg HTTPS_PROXY="${https_proxy}" \
  --build-arg NO_PROXY="${no_proxy}" \
  --build-arg http_proxy="${http_proxy}" \
  --build-arg https_proxy="${https_proxy}" \
  --build-arg no_proxy="${no_proxy}" \
  --tag tiller:%{IMAGE_TAG} \
  %{docker_build_dir}/tiller
mkdir -p %{docker_save_dir}/
docker save tiller:%{IMAGE_TAG} | xz -z -T2 > "%{docker_save_dir}/tiller:%{IMAGE_TAG}.tar"
docker rmi tiller:%{IMAGE_TAG}

%install
mkdir -p %{buildroot}/%{_caas_container_tar_path}
rsync -av %{docker_save_dir}/* %{buildroot}/%{_caas_container_tar_path}/

mkdir -p %{buildroot}/%{_roles_path}
rsync -av ansible/roles/helm %{buildroot}/%{_roles_path}/

install -D ansible/playbooks/helm.yaml %{buildroot}/%{_playbooks_path}/helm.yaml

install -D -m 0755 %{binary_build_dir}/helm %{buildroot}/usr/bin/helm

%files
%{_caas_container_tar_path}/tiller:%{IMAGE_TAG}.tar
%{_roles_path}/helm
%{_playbooks_path}/helm.yaml
/usr/bin/helm

%preun

%post
mkdir -p %{_postconfig_path}
ln -s %{_playbooks_path}/helm.yaml %{_postconfig_path}/

%postun
if [ $1 -eq 0 ]; then
  rm -f %{_postconfig_path}/helm.yaml
fi

%clean
rm -rf ${buildroot}

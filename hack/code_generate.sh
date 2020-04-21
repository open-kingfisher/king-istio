#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
#CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate the code with:
#Usage: $(basename "$0") <generators> <output-package> <apis-package> <groups-versions> ...
#
#  <generators>        the generators comma separated to run (deepcopy,defaulter,client,lister,informer) or "all".
#  <output-package>    the output package name (e.g. github.com/example/project/pkg/generated).
#  <apis-package>      the external types dir (e.g. github.com/example/api or github.com/example/project/pkg/apis).
#  <groups-versions>   the groups and their versions in the format "groupA:v1,v2 groupB:v1 groupC:v2", relative
#                      to <api-package>.
#  ...                 arbitrary flags passed to all generator binaries.
#
#
#Examples:
#  $(basename "$0") all             github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1"
#  $(basename "$0") deepcopy,client github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1"

# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.

# 关闭module模式，否则会出现panic: runtime error: index out of range问题
export GO111MODULE=off
# 如果没有generate-groups.sh　请使用　go get k8s.io/code-generator
# kingfisher/king-operator/pkg/client生成文件的路径，基于--output-base目录下面的
# kingfisher/king-operator/pkg/apis 定义的CRD接口文件路径
# 资源组的版本
# kingfisher/king-istio/pkg/apis 来自https://github.com/istio/client-go/tree/master/pkg/apis，因为里面已经存在deepcopy
# 所以脚本生成的时候不在生成deepcopy，而且types.gen.go 有关List 需要添加　+k8s:deepcopy-gen=true 标签才可以生成deepcopy，所以直接使用他的deepcopy
bash "${GOPATH}"/src/k8s.io/code-generator/generate-groups.sh "client,informer,lister" \
  kingfisher/king-istio/pkg/client kingfisher/king-istio/pkg/apis \
  "networking:v1alpha3 authentication:v1alpha1 rbac:v1alpha1 security:v1beta1 config:v1alpha2" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

# To use your own boilerplate text append:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt

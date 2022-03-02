#!/usr/bin/env bash
# Copyright 2022 DataStax
#
#   Licensed under the Apache License, Version 2.0 (the « License »);
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an « AS IS » BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

# script/install-astra.sh: dynamically install the correct binary according to the platform

EXE=astra
OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
fi

VERSION=$(curl --silent "https://api.github.com/repos/datastax-labs/astra-cli/releases/latest" |  grep tag_name | sed -nr 's/"tag_name": "(.+)",/\1/p'  | xargs)
VERSION_SHORT=${VERSION:1}

echo "installing $OS $ARCH $VERSION"
ARC_FOLDER=$EXE-cli_${VERSION_SHORT}_${OS}_${ARCH}
ARC=$(echo "${ARC_FOLDER}.tar.gz")

url=https://github.com/datastax-labs/astra-cli/releases/download/$VERSION/$ARC
curl -o $ARC -L $url
mkdir -p $ARC_FOLDER
tar zxvf $ARC  -C $ARC_FOLDER
sudo mv $ARC_FOLDER/$EXE /usr/local/bin/$EXE

rm -fr $ARC
rm -fr $ARC_FOLDER

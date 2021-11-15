#!/usr/bin/env bash
# Copyright 2021 Ryan Svihla
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

VERSION=$(curl --silent "https://api.github.com/repos/rsds143/astra-cli/releases/latest" |  grep tag_name | sed -nr 's/"tag_name": "(.+)",/\1/p'  | xargs)
VERSION_SHORT=${VERSION:1}

echo "installing $OS $ARCH $VERSION"
ARC_FOLDER=$EXE-${VERSION_SHORT}_${OS}_${ARCH}
ARC=$(echo "${ARC_FOLDER}.tar.gz")
if [ "$OS" = "darwin" ]; then
    ARC_FOLDER="$EXE-osx-$ARCH-signed"
    ARC=$(echo "${ARC_FOLDER}.zip")
    UNARC="unzip $ARC"
fi
url=https://github.com/rsds143/astra-cli/releases/download/$VERSION/$ARC
curl -O -L $url
if [ "$OS" = "darwin" ]; then
    ls
    cat astra-cli-osx-x86_64-signed.zip
    unzip $ARC
    sudo mv $EXE /usr/local/bin/$EXE
else 
    tar zxvf $ARC
    sudo mv $ARC_FOLDER/$EXE /usr/local/bin/$EXE
fi
rm -fr $ARC

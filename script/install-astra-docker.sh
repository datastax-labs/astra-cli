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

# script/install-astra-docker.sh: create a bash script that wraps the astra-cli docker image

sudo echo '#!/usr/bin/env bash' > /usr/local/bin/astra.sh
sudo echo 'docker run -it -v $HOME:/root ghcr.io/rsds143/astra-cli /astra "$@"' >> /usr/local/bin/astra.sh
sudo chmod +x /usr/local/bin/astra.sh

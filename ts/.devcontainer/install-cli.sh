#!/bin/bash

curl https://raw.githubusercontent.com/findy-network/findy-agent-cli/HEAD/install.sh > install.sh
chmod a+x install.sh
sudo ./install.sh -b /bin

asdf direnv setup --shell bash --version latest
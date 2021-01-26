#!/bin/bash

# Grab the binaries attached as release assets
curl -s https://api.github.com/repos/getcouragenow/main/releases/latest | jq -r ".assets[] | .browser_download_url" | xargs wget

# WORKING EXAMPLE :
curl -s https://api.github.com/repos/NFSTools/Binary/releases/latest | jq -r ".assets[] | .browser_download_url"  | xargs wget
zipfilename="$(curl -s https://api.github.com/repos/NFSTools/Binary/releases/latest | jq -r '.assets[].name')"
unzip -j $zipfilename -d binaries

# Get the code as tar.gz, decompress..
curl -L "$(curl -s https://api.github.com/repos/applinh-getcouragenow/main/releases/latest | jq -r ".tarball_url")" > main.tar.gz 
chmod 700 main.tar.gz
mkdir main
tar xzvf main.tar.gz -C main --strip-components=1 
# .. then maybe cd inside it
# run make command
# then delete everything keep only the binaries
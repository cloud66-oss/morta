#!/bin/bash

if [ -z "$2" ]
  then
    version=$(git describe --tags --always)
  else
    version=$2
fi

if [ -z "$1" ]
  then
    echo "No channel supplied"
    exit 1
fi

channel=$1
force="false"

if [[ $2 == "--force" ]]
  then
    force="true"
fi

echo "Building $channel/$version"
echo

rm build/*

curl -s http://downloads.cloud66.com.s3.amazonaws.com/morta/versions.json | jq '.versions |= map(if (.channel == "'$channel'") then .version = "'$version'" else . end) | .versions |= map(if (.channel == "'$channel'") then .force = '$force' else . end)' > build/versions.json
echo "Current Versions"
cat build/versions.json | jq -r '.versions | map([.channel, .version] | join(": ")) | .[]'
echo

gox -ldflags "-X github.com/cloud66-oss/morta/utils.Version=$version -X github.com/cloud66-oss/morta/utils.Channel=$channel" -os="darwin linux" -arch="amd64" -output "build/{{.OS}}_{{.Arch}}_$version"

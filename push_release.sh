#!/bin/bash

set -e

err_report() {
	echo "Error on line $1"
}

trap 'err_report $LINENO' ERR

echo "building release artifacts"
PROJECTDIR=${TRAVIS_HOME}/gopath/src/github.com/${TRAVIS_REPO_SLUG}
echo ${PROJECTDIR}

# prepare artifacts with goreleaser, results will be stored to the ./dist directory
curl -sL https://git.io/goreleaser | bash -s -- --snapshot --skip-publish --rm-dist

# build an AMI image with packer
curl -L -o packer.zip https://releases.hashicorp.com/packer/1.4.1/packer_1.4.1_linux_amd64.zip
unzip packer.zip
./packer build -force build/packer.json

# if a tag has alpha or beta in the name, it will be released as a pre-release.
# if a tag does not have alpha or beta, it is pushed as a full release.
case "${TAG}" in
	*alpha* )  echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG"  dist/;;
	*beta* )    echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG"   dist/;;
	*)echo "Releasing version: ${TAG}, as latest release."
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "latest release" --debug "$TAG"   dist/;;
esac

# Check for errors
if [ $? -eq 0 ]; then
	echo "Release pushed"
else
	echo "Push to releases failed"
	exit 1
fi

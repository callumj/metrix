#!/bin/bash
pattern='Size: ([0-9]+)'

mkdir -p tmp
zip tmp/assets.zip assets/*

stat_res=""
arch=`uname`
if [[ $arch == "Darwin" ]]; then
  stat_res=`stat -x tmp/assets.zip`
else
  stat_res=`stat tmp/assets.zip`
fi

fileLength=false
if [[ "$stat_res" =~ $pattern ]]; then
  fileLength=${BASH_REMATCH[1]}
fi

if [[ fileLength == false ]]; then
  echo "Cannot figure out size"
  exit -1
else
  echo "Size of archive: ${fileLength}"
fi

VERSION=`cat VERSION`

rm -r -f tmp/go
rm -r -f builds/

# vet the source (capture errors because the current version does not use exit statuses currently)
echo "Vetting..."
VET=`go tool vet . 2>&1 >/dev/null`

cur=`pwd` 

if ! [ -n "$VET" ]
then
  echo "All good"
  mkdir tmp/go
  mkdir tmp/go/src tmp/go/bin tmp/go/pkg
  mkdir -p tmp/go/src/github.com/callumj/metrix
  cp -R app assets handlers metric_core resource_bundle shared main.go tmp/go/src/github.com/callumj/metrix/
  mkdir -p builds/${VERSION}/darwin_386 builds/${VERSION}/darwin_amd64 builds/${VERSION}/linux_386 builds/${VERSION}/linux_amd64 
  GOPATH="${cur}/tmp/go"
  echo "Getting"
  GOPATH="${cur}/tmp/go" go get -d .
  echo "Starting build"

  GOPATH="${cur}/tmp/go" GOOS=darwin GOARCH=386 go build -o builds/${VERSION}/darwin_386/metrix

  GOPATH="${cur}/tmp/go" GOARCH=amd64 GOOS=darwin go build -o builds/${VERSION}/darwin_amd64/metrix

  GOPATH="${cur}/tmp/go" GOOS=linux GOARCH=amd64 go build -o builds/${VERSION}/linux_amd64/metrix

  GOPATH="${cur}/tmp/go" GOOS=linux GOARCH=386 go build -o builds/${VERSION}/linux_386/metrix
else
  echo "$VET"
  exit -1
fi

# rewrite the binaries

FILES=builds/${VERSION}/*/metrix
for f in $FILES
do
  cat tmp/assets.zip >> ${f}
  echo -n "ArchiveLength:${fileLength}" >> ${f}
  str="/metrix"
  repl=""
  path=${f/$str/$repl}
  tar  -C ${path} -cvzf "${f}.tgz" metrix
  rm ${f}
done
#!/bin/bash
pattern='Size: ([0-9]+)'

mkdir -p tmp
zip tmp/assets.zip assets/*

stat_res=`stat -x tmp/assets.zip`

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

# vet the source (capture errors because the current version does not use exit statuses currently)
VET=`go tool vet . 2>&1 >/dev/null`

rm -r -f builds/

if ! [ -n "$VET" ]
then
  echo "All good"
  goxc -os "linux darwin" -pv ${VERSION} -d builds xc copy-resources
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
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

/bin/bash --login -c "ginst"

cat tmp/assets.zip >> `which metrix`
echo -n "ArchiveLength:${fileLength}" >> `which metrix`

metrix config.yml
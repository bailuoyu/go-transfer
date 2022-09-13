#!/bin/bash
branch=""
if [ ! -n "$1" ]; then
  branch=$(git rev-parse --abbrev-ref HEAD)
  echo "Now Origin Branch: ${branch}"
else
  branch=$1
fi
git reset --hard origin/${branch}
echo "确认分支信息开始下一步"
read answer
git pull

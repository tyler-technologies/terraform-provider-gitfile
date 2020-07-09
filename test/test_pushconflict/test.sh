#!/bin/bash
printf "\nRUNNING TEST_PUSHCONFLICT\n"
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SCRIPTPATH/create_conflict.sh
set -ex
rm -rf example.git checkout terraform.tfstate terraform.tfstate.backup
mkdir example.git
cd example.git
git init
touch .exists
git add .exists
git commit -m"Initial commit"
git checkout -b move_HEAD
cd ..
terraform init
terraform apply -auto-approve  & create_conflict
cd checkout
git log | grep 'Created by terraform gitfile_commit'
git fetch
git log origin/master | grep 'Created by terraform gitfile_commit'
if [ ! -f terraform ]; then
    exit 1
fi
sleep 2
cd ..
terraform destroy -auto-approve
sleep 2
if [ -d checkout ]; then
    exit 1
fi

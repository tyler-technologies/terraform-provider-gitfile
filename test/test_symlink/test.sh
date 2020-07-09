#!/bin/bash
printf "\nRUNNING TEST_SYMLINK\n"
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
terraform apply -auto-approve
terraform apply -auto-approve
cd checkout
git fetch
# We did do a commit
git log origin/master | grep 'Created by terraform gitfile_commit'
if [ ! -L terraform ]; then
    exit 1
fi
if [ "$(readlink terraform)" != "/etc/passwd" ]; then
    exit 1
fi

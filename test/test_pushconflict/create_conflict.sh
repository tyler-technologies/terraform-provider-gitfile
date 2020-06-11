#!/bin/bash

git_path=./example.git

create_commit() {
  git -C ${git_path} checkout master
  echo "conflict changes ${1}" > example.git/terraform
  git -C ${git_path} add terraform
  git -C ${git_path} commit -m "test conflict $1"
  git -C ${git_path} checkout move_HEAD
}

create_conflict() {
  for i in $(seq 1 10); do
    create_commit ${i}
    sleep 1
  done
}
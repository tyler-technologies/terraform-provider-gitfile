#!/bin/bash

git_path=./example.git

create_conflict() {
  git -C ${git_path} checkout master
  echo "conflict changes" > example.git/terraform
  git -C ${git_path} add terraform
  git -C ${git_path} commit -m "test conflict"
  git -C ${git_path} checkout move_HEAD
}
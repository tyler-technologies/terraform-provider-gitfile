provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "checkout" {}

resource "gitfile_file" "testfile" {
    path = "terraform"
    contents = "Terraform making commits"
}

resource "gitfile_file" "shizz" {
    path = "myfile"
    contents = "Terraform shizz"
}

resource "gitfile_commit" "commit" {
    commit_message = "Created by terraform gitfile_commit"
    # handles = ["${gitfile_file.testfile.id}"]
    handles = ["${gitfile_file.testfile.id}", "${gitfile_file.shizz.id}"]
}

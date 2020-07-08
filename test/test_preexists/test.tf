provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "test" {}

resource "gitfile_file" "test" {
    path = "terraform"
    contents = "preexisting_commits\n"
}
resource "gitfile_commit" "test" {
    commit_message = "Created by terraform gitfile_commit"
    handles = ["${gitfile_file.test.id}"]
}

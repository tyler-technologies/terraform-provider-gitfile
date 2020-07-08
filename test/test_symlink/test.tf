provider "gitfile" {
    repo_url = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_checkout" "test" {}

resource "gitfile_symlink" "test" {
    path = "terraform"
    target = "/etc/passwd"
}
resource "gitfile_commit" "test" {
    commit_message = "Created by terraform gitfile_commit"
    handles = [gitfile_symlink.test.id]
}


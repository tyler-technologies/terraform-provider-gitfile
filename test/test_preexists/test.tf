resource "gitfile_checkout" "test" {
    repo = "../example.git"
    branch = "master"
    path = "checkout"
}
resource "gitfile_file" "test" {
    checkout_dir = gitfile_checkout.test.path
    path = "terraform"
    contents = "preexisting_commits\n"
}
resource "gitfile_commit" "test" {
    checkout_dir = gitfile_checkout.test.path
    commit_message = "Created by terraform gitfile_commit"
    handles = ["${gitfile_file.test.id}"]
}

resource "gitfile_checkout" "checkout" {
    repo = "../example.git"
    branch = "master"
    path = "checkout"
}

resource "gitfile_file" "testfile" {
    checkout_dir = gitfile_checkout.checkout.path
    path = "terraform"
    contents = "Terraform making commits"
}

resource "gitfile_file" "shizz" {
    checkout_dir = gitfile_checkout.checkout.path
    path = "myfile"
    contents = "Terraform shizz"
}

resource "gitfile_commit" "commit" {
    checkout_dir = gitfile_checkout.checkout.path
    commit_message = "Created by terraform gitfile_commit"
    # handles = ["${gitfile_file.testfile.id}"]
    handles = ["${gitfile_file.testfile.id}", "${gitfile_file.shizz.id}"]
}

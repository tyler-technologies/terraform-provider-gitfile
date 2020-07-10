package simple

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-technologies/terraform-provider-gitfile/test/helpers"
)

func setup() {
	cleanup()
	example_dir := "example.git"

	os.MkdirAll(example_dir, 0755)
	helpers.GitCommand(example_dir, "init")
	os.Create(fmt.Sprintf("%s/.exists", example_dir))
	helpers.GitCommand(example_dir, "add", ".exists")
	helpers.GitCommand(example_dir, "commit", "-m", "Initial Commit")
	helpers.GitCommand(example_dir, "checkout", "-b", "move_HEAD")
}

func TestSimple(t *testing.T) {
	setup()

	o := &helpers.TerratestDefaultOptions
	terraform.Init(t, o)
	helpers.GeneratePlan(t, o, "plan.out")
	helpers.ApplyWithPlanFile(t, o, "plan.out")
	// terraform.InitAndApply(t, o)
	expected_commit_msg := "Created by terraform gitfile_commit"

	tests := []struct {
		output   string
		expected string
	}{
		{"gitfile_checkout_path", "checkout"},
		{"gitfile_commit_commit_message", expected_commit_msg},
	}

	for _, test := range tests {
		actual, _ := terraform.OutputE(t, o, test.output)
		assert.Equal(t, test.expected, actual, "terraform output '%s'", test.output)
	}
	found, err := helpers.GitLogContains("checkout", expected_commit_msg)
	assert.NoError(t, err)
	assert.True(t, found, fmt.Sprintf("checkout should have commit message '%s'", expected_commit_msg))
}

func cleanup() {
	os.RemoveAll("example.git")
	os.RemoveAll("checkout")
	os.RemoveAll("terraform.tfstate")
	os.RemoveAll("terraform.tfstate.backup")
	os.RemoveAll("plan.out")
}

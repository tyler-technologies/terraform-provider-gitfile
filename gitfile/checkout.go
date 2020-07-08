package gitfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func checkoutResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					i := strings.IndexRune(value, '/')
					if i == 0 {
						es = append(es, fmt.Errorf("Paths which begin with / not allowed in %q", k))
					}
					return
				},
			},
			"repo": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
				ForceNew: true, // FIXME
			},
			"head": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retry_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     10,
				Description: "Number of git commit retries",
			},

			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     5,
				Description: "Number of seconds between git commit retries",
			},
		},
		Create: CheckoutCreate,
		Read:   CheckoutRead,
		Update: nil,
		Delete: CheckoutDelete,
	}
}

func clone(dir, repo, branch string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// May already be checked out from another project
	if _, err := os.Stat(fmt.Sprintf("%s/.git", dir)); err != nil {
		if _, err := gitCommand(dir, "clone", "-b", branch, "--", repo, "."); err != nil {
			return err
		}
	}
	return nil
}

func read(d *schema.ResourceData) error {
	checkout_dir := d.Id()
	var repo string
	var branch string
	var head string

	if out, err := gitCommand(checkout_dir, "config", "--get", "remote.origin.url"); err != nil {
		return err
	} else {
		repo = strings.TrimRight(string(out), "\n")
	}
	if out, err := gitCommand(checkout_dir, "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return err
	} else {
		branch = strings.TrimRight(string(out), "\n")
	}

	if _, err := gitCommand(checkout_dir, "pull", "--ff-only", "origin"); err != nil {
		return err
	}

	if out, err := gitCommand(checkout_dir, "rev-parse", "HEAD"); err != nil {
		return err
	} else {
		head = strings.TrimRight(string(out), "\n")
	}
	d.Set("path", checkout_dir)
	d.Set("repo", repo)
	d.Set("branch", branch)
	d.Set("head", head)
	return nil
}

func CheckoutCreate(d *schema.ResourceData, meta interface{}) error {
	checkout_dir := d.Get("path").(string)
	repo := d.Get("repo").(string)
	branch := d.Get("branch").(string)

	err := clone(checkout_dir, repo, branch)
	if err != nil {
		return err
	}

	d.SetId(checkout_dir)
	return read(d)
}

func CheckoutRead(d *schema.ResourceData, meta interface{}) error {
	checkout_dir := d.Id()
	repo := d.Get("repo").(string)
	branch := d.Get("branch").(string)

	if _, err := os.Stat(checkout_dir); err != nil {
		err = clone(checkout_dir, repo, branch)
		if err != nil {
			return err
		}
	}
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	read(d)

	return nil
}

func CheckoutDelete(d *schema.ResourceData, meta interface{}) error {
	checkout_dir := d.Id()
	retry_count := d.Get("retry_count").(int)
	retry_interval := d.Get("retry_interval").(int)
	expected_repo := d.Get("repo").(string)
	expected_branch := d.Get("branch").(string)
	expected_head := d.Get("head").(string)

	if _, err := os.Stat(checkout_dir); err != nil {
		err = clone(checkout_dir, expected_repo, expected_branch)
		if err != nil {
			return err
		}
	}

	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)

	// sanity check
	var repo string
	var branch string
	var head string

	if out, err := gitCommand(checkout_dir, "config", "--get", "remote.origin.url"); err != nil {
		return err
	} else {
		repo = strings.TrimRight(string(out), "\n")
	}
	if out, err := gitCommand(checkout_dir, "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return err
	} else {
		branch = strings.TrimRight(string(out), "\n")
	}

	if _, err := gitCommand(checkout_dir, "pull", "--ff-only", "origin"); err != nil {
		return err
	}

	if out, err := gitCommand(checkout_dir, "rev-parse", "HEAD"); err != nil {
		return err
	} else {
		head = strings.TrimRight(string(out), "\n")
	}

	if expected_repo != repo {
		return fmt.Errorf("expected repo to be %s, was %s", expected_repo, repo)
	}
	if expected_branch != branch {
		return fmt.Errorf("expected branch to be %s, was %s", expected_branch, branch)
	}
	if expected_head != head {
		return fmt.Errorf("expected head to be %s, was %s", expected_head, head)
	}

	if err := commit(checkout_dir, "Removed by Terraform", ""); err != nil {
		return errwrap.Wrapf("push error: {{err}}", err)
	}

	if err := push(checkout_dir, 0, retry_count, retry_interval); err != nil {
		return err
	}

	// actually delete
	if err := os.RemoveAll(checkout_dir); err != nil {
		return err
	}

	return nil
}

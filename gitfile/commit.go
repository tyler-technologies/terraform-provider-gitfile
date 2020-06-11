package gitfile

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const CommitBodyHeader string = "The following files are managed by terraform:"

func commitResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"commit_message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Created by terraform gitfile_commit",
				ForceNew: true,
			},
			"checkout_dir": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"handles": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      hashString,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		Create: CommitCreate,
		Read:   CommitRead,
		Delete: CommitDelete,
		Exists: CommitExists,
	}
}

func CommitCreate(d *schema.ResourceData, meta interface{}) error {
	checkout_dir := d.Get("checkout_dir").(string)
	retry_count := d.Get("retry_count").(int)
	retry_interval := d.Get("retry_interval").(int)
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)

	handles := d.Get("handles").(*schema.Set)
	commit_message := d.Get("commit_message").(string)
	filepaths := []string{}
	for _, handle := range handles.List() {
		filepaths = append(filepaths, parseHandle(handle.(string)).path)
	}

	var sha string

	if err := stash(checkout_dir); err != nil {
		return err
	}

	if err := pull(checkout_dir); err != nil {
		return err
	}

	if err := applyStash(checkout_dir); err != nil {
		return err
	}

	commit_body := fmt.Sprintf("%s\n%s", CommitBodyHeader, strings.Join(filepaths, "\n"))
	if _, err := gitCommand(checkout_dir, flatten("commit", "-m", commit_message, "-m", commit_body, "--allow-empty")...); err != nil {
		return err
	}

	if err := push(checkout_dir, 0, retry_count, retry_interval); err != nil {
		return err
	}

	if out, err := gitCommand(checkout_dir, "rev-parse", "HEAD"); err != nil {
		return err
	} else {
		sha = strings.TrimRight(string(out), "\n")
	}

	d.SetId(fmt.Sprintf("%s %s", sha, checkout_dir))
	return nil
}

func CommitRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func CommitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	checkoutDir := d.Get("checkout_dir").(string)
	lockCheckout(checkoutDir)
	defer unlockCheckout(checkoutDir)
	commitId := strings.Split(d.Id(), " ")[0]

	_, err := gitCommand(checkoutDir, flatten("show", commitId)...)

	if err != nil {
		return false, nil
	} else {
		return true, nil
	}

}

func CommitDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func push(checkout_dir string, count int8, retry_count, retry_interval int) error {
	if _, err := gitCommand(checkout_dir, "push", "origin", "HEAD"); err != nil {
		if int(count) >= retry_count {
			return err
		}

		time.Sleep(time.Duration(retry_interval) * time.Second)
		count++

		if err := resetCommit(checkout_dir); err != nil {
			return errwrap.Wrapf("push error: {{err}}", err)
		}

		if err := stash(checkout_dir); err != nil {
			return errwrap.Wrapf("push error: {{err}}", err)
		}

		if err := pull(checkout_dir); err != nil {
			return errwrap.Wrapf("push error: {{err}}", err)
		}

		if err := applyStash(checkout_dir); err != nil {
			return errwrap.Wrapf("push error: {{err}}", err)
		}

		return push(checkout_dir, count, retry_count, retry_interval)
	}
	return nil
}

func stash(checkout_dir string) error {
	if _, err := gitCommand(checkout_dir, "stash"); err != nil {
		return err
	}
	return nil
}

func pull(checkout_dir string) error {
	if _, err := gitCommand(checkout_dir, "pull"); err != nil {
		return err
	}
	return nil
}

func resetCommit(checkout_dir string) error {
	if _, err := gitCommand(checkout_dir, "reset", "--soft", "HEAD~1"); err != nil {
		return err
	}
	return nil
}

func applyStash(checkout_dir string) error {
	if _, err := gitCommand(checkout_dir, "stash", "show", "stash@{0}"); err != nil {
		return nil
	}

	if _, err := gitCommand(checkout_dir, "checkout", "stash", "--", "."); err != nil {
		return err
	}
	return nil
}

package gitfile

import (
	"os"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func symlinkResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: symlinkCreateUpdate,
		Read:   symlinkRead,
		Delete: symlinkDelete,
		Exists: symlinkExists,
	}
}

func symlinkCreateUpdate(d *schema.ResourceData, m interface{}) error {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)

	filepath := d.Get("path").(string)
	target := d.Get("target").(string)

	if err := os.Remove(path.Join(checkout_dir, filepath)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(path.Dir(path.Join(checkout_dir, filepath)), 0755); err != nil {
		return err
	}
	if err := os.Symlink(target, path.Join(checkout_dir, filepath)); err != nil {
		return err
	}

	if _, err := gitCommand(checkout_dir, "add", "--", filepath); err != nil {
		return err
	}

	hand := handle{
		kind: "symlink",
		hash: hashString(target),
		path: filepath,
	}

	d.SetId(hand.String())
	return nil
}

func symlinkRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func symlinkExists(d *schema.ResourceData, m interface{}) (bool, error) {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	filepath := d.Get("path").(string)
	var target string
	var err error
	if target, err = os.Readlink(path.Join(checkout_dir, filepath)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	if target == d.Get("target").(string) {
		return true, nil
	} else {
		return false, nil
	}
}

func symlinkDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

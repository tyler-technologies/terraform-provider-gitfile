package gitfile

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func fileResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contents": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: fileCreateUpdate,
		Read:   fileRead,
		Delete: fileDelete,
		Exists: fileExists,
	}
}

func fileCreateUpdate(d *schema.ResourceData, m interface{}) error {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)

	filepath := d.Get("path").(string)
	contents := d.Get("contents").(string)

	if err := os.MkdirAll(path.Dir(path.Join(checkout_dir, filepath)), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(checkout_dir, filepath), []byte(contents), 0666); err != nil {
		return err
	}

	if _, err := gitCommand(checkout_dir, "add", "--", filepath); err != nil {
		return err
	}

	hand := handle{
		kind: "file",
		hash: hashString(contents),
		path: filepath,
	}

	d.SetId(hand.String())
	return nil
}

func fileRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func fileExists(d *schema.ResourceData, m interface{}) (bool, error) {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	filepath := d.Get("path").(string)

	var out []byte
	var err error
	if out, err = ioutil.ReadFile(path.Join(checkout_dir, filepath)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if string(out) == d.Get("contents").(string) {
		return true, nil
	} else {
		return false, nil
	}
}

func fileDelete(d *schema.ResourceData, m interface{}) error {
	checkout_dir := m.(*GitFileConfig).Path
	lockCheckout(checkout_dir)
	defer unlockCheckout(checkout_dir)
	filepath := d.Get("path").(string)
	file := path.Join(checkout_dir, filepath)
	if err := os.Remove(file); err != nil {
		return err
	}
	if _, err := gitCommand(checkout_dir, "rm", "--", filepath); err != nil {
		return err
	}
	return nil
}

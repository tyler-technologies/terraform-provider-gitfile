package gitfile

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"commit_retry_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Number of git commit retries",
			},

			"commit_retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Number of seconds between git commit retries",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gitfile_checkout": checkoutResource(),
			"gitfile_file":     fileResource(),
			"gitfile_symlink":  symlinkResource(),
			"gitfile_commit":   commitResource(),
		},
		ConfigureFunc: gitfileConfigure,
	}
}

func gitfileConfigure(data *schema.ResourceData) (interface{}, error) {
	config := &gitfileConfig{
		CommitRetryCount:    data.Get("commit_retry_count").(int8),
		CommitRetryInterval: data.Get("commit_retry_interval").(int8),
	}
	return config, nil
}

type gitfileConfig struct {
	CommitRetryCount    int8
	CommitRetryInterval int8
}

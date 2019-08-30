package hiera

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hiera.yaml",
			},
			"bin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lookup",
			},
			"scope": {
				Type:     schema.TypeMap,
				Default:  map[string]interface{}{},
				Optional: true,
			},
			"merge": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "first",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hiera":       dataSourceHiera(),
			"hiera_array": dataSourceHieraArray(),
			"hiera_hash":  dataSourceHieraHash(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	return NewHiera(
		data.Get("bin").(string),
		data.Get("config").(string),
		data.Get("scope").(map[string]interface{}),
		data.Get("merge").(string),
	), nil
}

package hiera

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHiera() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHieraRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHieraRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera value")

	keyName := d.Get("key").(string)

	config := meta.(Config)
	v, err := config.Value(keyName)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(keyName)
	d.Set("value", v)

	return nil
}

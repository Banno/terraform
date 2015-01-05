package logical

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogicalMutate() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalMutateCreate,
		Read:   resourceLogicalMutateRead,
		Update: resourceLogicalMutateUpdate,
		Delete: resourceLogicalMutateDelete,

		Schema: map[string]*schema.Schema{
			// look into using the connection object instead.
			// create the terraform-mutate.sh in the destination folder
			"ssh": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_file": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"source": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"env": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"exec": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"force": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},
		},
	}
}

func resourceLogicalMutateCreate(d *schema.ResourceData, meta interface{}) error {

	return resourceLogicalMutateUpdate(d, meta)
}

func resourceLogicalMutateRead(d *schema.ResourceData, meta interface{}) error {
	// Make sure the unique ID is up to date and provide force functionality.

	d.SetId(fmt.Sprintf("%s -- %s", d.Get("ssh.0.ip").(string), d.Get("exec").(string)))
	d.Set("force", "false")

	return nil
}

func resourceLogicalMutateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("======> ResourceData  <========\n%#v\n", d)

	// get access to the connection object and read some remote data

	// copy a file
	// copy a folder

	// render the env variables into a template with the exec command and make it executable
	// execute the script

	return resourceLogicalMutateRead(d, meta)
}

func resourceLogicalMutateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

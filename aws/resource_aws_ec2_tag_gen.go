// Code generated by internal/tagresource/generator/main.go; DO NOT EDIT.

package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/aws/internal/tags"
	tftags "github.com/hashicorp/terraform-provider-aws/aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagCreate,
		Read:   resourceTagRead,
		Update: resourceTagUpdate,
		Delete: resourceTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	identifier := d.Get("resource_id").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	if err := tftags.Ec2CreateTags(conn, identifier, map[string]string{key: value}); err != nil {
		return fmt.Errorf("error creating %s resource (%s) tag (%s): %w", ec2.ServiceID, identifier, key, err)
	}

	d.SetId(tftags.SetResourceID(identifier, key))

	return resourceTagRead(d, meta)
}

func resourceTagRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	identifier, key, err := tftags.GetResourceID(d.Id())

	if err != nil {
		return err
	}

	value, err := tftags.Ec2GetTag(conn, identifier, key)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] %s resource (%s) tag (%s) not found, removing from state", ec2.ServiceID, identifier, key)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading %s resource (%s) tag (%s): %w", ec2.ServiceID, identifier, key, err)
	}

	d.Set("resource_id", identifier)
	d.Set("key", key)
	d.Set("value", value)

	return nil
}

func resourceTagUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	identifier, key, err := tftags.GetResourceID(d.Id())

	if err != nil {
		return err
	}

	if err := tftags.Ec2UpdateTags(conn, identifier, nil, map[string]string{key: d.Get("value").(string)}); err != nil {
		return fmt.Errorf("error updating %s resource (%s) tag (%s): %w", ec2.ServiceID, identifier, key, err)
	}

	return resourceTagRead(d, meta)
}

func resourceTagDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	identifier, key, err := tftags.GetResourceID(d.Id())

	if err != nil {
		return err
	}

	if err := tftags.Ec2UpdateTags(conn, identifier, map[string]string{key: d.Get("value").(string)}, nil); err != nil {
		return fmt.Errorf("error deleting %s resource (%s) tag (%s): %w", ec2.ServiceID, identifier, key, err)
	}

	return nil
}

package provider

import (
	"fmt"

	"github.com/sergo-techhub/updown"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func statusPageResource() *schema.Resource {
	return &schema.Resource{
		Description: "`updown_status_page` defines a status page",

		Create: statusPageCreate,
		Read:   statusPageRead,
		Delete: statusPageDelete,
		Update: statusPageUpdate,
		Exists: statusPageExists,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"checks": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of checks to show in the page (order is respected).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the status page.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description text (displayed below the name, supports newlines and links).",
			},
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Page visibility: 'public', 'protected', or 'private'.",
				Default:     "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public", "protected", "private",
				}, false),
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access key for protected pages (defaults to a random string if unset).",
				Sensitive:   true,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the status page.",
			},
		},
	}
}

func constructStatusPagePayload(d *schema.ResourceData) updown.StatusPageItem {
	payload := updown.StatusPageItem{}

	if v, ok := d.GetOk("checks"); ok {
		interfaceSlice := v.([]interface{})
		var stringSlice []string
		for _, s := range interfaceSlice {
			stringSlice = append(stringSlice, s.(string))
		}
		payload.Checks = stringSlice
	}

	if v, ok := d.GetOk("name"); ok {
		payload.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		payload.Description = v.(string)
	}

	if v, ok := d.GetOk("visibility"); ok {
		payload.Visibility = v.(string)
	}

	if v, ok := d.GetOk("access_key"); ok {
		payload.AccessKey = v.(string)
	}

	return payload
}

func statusPageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)

	statusPage, _, err := client.StatusPage.Add(constructStatusPagePayload(d))
	if err != nil {
		return fmt.Errorf("creating status page with the API: %w", err)
	}

	d.SetId(statusPage.Token)

	return statusPageRead(d, meta)
}

func statusPageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)
	statusPage, _, err := client.StatusPage.Get(d.Id())

	if err != nil {
		return fmt.Errorf("reading status page from the API: %w", err)
	}

	for k, v := range map[string]interface{}{
		"checks":      statusPage.Checks,
		"name":        statusPage.Name,
		"description": statusPage.Description,
		"visibility":  statusPage.Visibility,
		"access_key":  statusPage.AccessKey,
		"url":         statusPage.URL,
	} {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func statusPageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)

	_, _, err := client.StatusPage.Update(d.Id(), constructStatusPagePayload(d))
	if err != nil {
		return fmt.Errorf("updating status page with the API: %w", err)
	}

	return statusPageRead(d, meta)
}

func statusPageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)
	deleted, _, err := client.StatusPage.Remove(d.Id())

	if err != nil {
		return fmt.Errorf("removing status page from the API: %w", err)
	}

	if !deleted {
		return fmt.Errorf("status page couldn't be deleted")
	}

	return nil
}

func statusPageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := statusPageRead(d, meta)
	return err == nil, err
}

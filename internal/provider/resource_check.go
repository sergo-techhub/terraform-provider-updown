package provider

import (
	"fmt"
	"strings"

	"github.com/sergo-techhub/updown"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func checkResource() *schema.Resource {
	return &schema.Resource{
		Description: "`updown_check` defines a check",

		Create: checkCreate,
		Read:   checkRead,
		Delete: checkDelete,
		Update: checkUpdate,
		Exists: checkExists,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL you want to monitor.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Interval in seconds (15, 30, 60, 120, 300, 600, 1800 or 3600).",
				Default:     60,
			},
			"apdex_t": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "APDEX threshold in seconds (0.125, 0.25, 0.5, 1.0 or 2.0).",
				Default:     0.5,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the check enabled (true or false).",
				Default:     true,
			},
			"published": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Shall the status page be public (true or false).",
				Default:     false,
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human readable name.",
			},
			"string_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search for this string in the page.",
			},
			"mute_until": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mute notifications until given time, accepts a time, 'recovery' or 'forever'.",
			},
			"disabled_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Disabled monitoring locations. It's a lsit of abbreviated location names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recipients": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Selected alert recipients. It's an array of recipient IDs you can get from the recipients API.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The HTTP headers you want in requests.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of check (http, https, icmp, tcp, tcps). Inferred from URL scheme if not specified.",
				ValidateFunc: validation.StringInSlice([]string{
					"http", "https", "icmp", "tcp", "tcps",
				}, false),
			},
			"http_verb": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP method (GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS). Only for http/https checks.",
				Default:     "GET",
				ValidateFunc: validation.StringInSlice([]string{
					"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
				}, false),
			},
			"http_body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Request body for POST/PUT/PATCH requests. Only for http/https checks.",
			},
		},
	}
}

func constructCheckPayload(d *schema.ResourceData) updown.CheckItem {
	payload := updown.CheckItem{}

	if v, ok := d.GetOk("url"); ok {
		payload.URL = v.(string)
	}

	if v, ok := d.GetOk("period"); ok {
		payload.Period = v.(int)
	}

	if v, ok := d.GetOk("apdex_t"); ok {
		payload.Apdex = v.(float64)
	}

	if v, ok := d.GetOk("enabled"); ok {
		payload.Enabled = v.(bool)
	}

	if v, ok := d.GetOk("published"); ok {
		payload.Published = v.(bool)
	}

	if v, ok := d.GetOk("alias"); ok {
		payload.Alias = v.(string)
	}

	if v, ok := d.GetOk("string_match"); ok {
		payload.StringMatch = v.(string)
	}

	if v, ok := d.GetOk("mute_until"); ok {
		payload.MuteUntil = v.(string)
	}

	if v, ok := d.GetOk("disabled_locations"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var stringSlice []string
		for _, s := range interfaceSlice {
			stringSlice = append(stringSlice, s.(string))
		}
		payload.DisabledLocations = stringSlice
	}

	if v, ok := d.GetOk("recipients"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var stringSlice []string
		for _, s := range interfaceSlice {
			stringSlice = append(stringSlice, s.(string))
		}
		payload.RecipientIDs = stringSlice
	}

	if m, ok := d.GetOk("custom_headers"); ok {
		payload.CustomHeaders = map[string]string{}
		for k, v := range m.(map[string]interface{}) {
			payload.CustomHeaders[k] = v.(string)
		}
	}

	checkType := ""
	if v, ok := d.GetOk("type"); ok {
		checkType = v.(string)
		payload.Type = checkType
	}

	// Only set http_verb and http_body for HTTP/HTTPS checks
	// Don't send http_verb if it's the default GET - API defaults to GET
	isHttpCheck := checkType == "" || checkType == "http" || checkType == "https"
	if isHttpCheck {
		httpVerb := d.Get("http_verb").(string)
		if httpVerb != "" && httpVerb != "GET" {
			payload.HttpVerb = httpVerb
		}

		if v, ok := d.GetOk("http_body"); ok {
			payload.HttpBody = v.(string)
		}
	}

	return payload
}

func checkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)

	check, _, err := client.Check.Add(constructCheckPayload(d))
	if err != nil {
		return fmt.Errorf("creating check with the API: %w", err)
	}

	d.SetId(check.Token)

	return checkRead(d, meta)
}

func checkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)
	check, _, err := client.Check.Get(d.Id())

	if err != nil {
		return fmt.Errorf("reading check from the API: %w", err)
	}

	// Normalize URL by stripping protocol prefix for non-HTTP checks
	// The API returns URLs like "icmp://192.168.1.1" but we store just "192.168.1.1"
	normalizedURL := check.URL
	isHttpCheck := check.Type == "" || check.Type == "http" || check.Type == "https"

	if check.Type == "icmp" {
		normalizedURL = strings.TrimPrefix(normalizedURL, "icmp://")
	}

	// For non-HTTP checks, set http_verb to match schema default to prevent drift
	httpVerb := check.HttpVerb
	httpBody := check.HttpBody
	if !isHttpCheck {
		httpVerb = "GET" // Match schema default
		httpBody = ""
	}

	for k, v := range map[string]interface{}{
		"url":                normalizedURL,
		"period":             check.Period,
		"apdex_t":            check.Apdex,
		"enabled":            check.Enabled,
		"published":          check.Published,
		"alias":              check.Alias,
		"string_match":       check.StringMatch,
		"mute_until":         check.MuteUntil,
		"disabled_locations": check.DisabledLocations,
		"recipients":         check.RecipientIDs,
		"custom_headers":     check.CustomHeaders,
		"type":               check.Type,
		"http_verb":          httpVerb,
		"http_body":          httpBody,
	} {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func checkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)

	_, _, err := client.Check.Update(d.Id(), constructCheckPayload(d))
	if err != nil {
		return fmt.Errorf("updating check with the API: %w", err)
	}

	return nil
}

func checkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*updown.Client)
	checkDeleted, _, err := client.Check.Remove(d.Id())

	if err != nil {
		return fmt.Errorf("removing check from the API: %w", err)
	}

	if !checkDeleted {
		return fmt.Errorf("check couldn't be deleted")
	}

	return nil
}

func checkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := checkRead(d, meta)
	return err == nil, err
}

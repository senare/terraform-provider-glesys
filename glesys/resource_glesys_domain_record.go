package glesys

import (
	"context"
	"fmt"
	"strconv"

	"github.com/glesys/glesys-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysDNSDomainRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysDNSDomainRecordCreate,
		Read:   resourceGlesysDNSDomainRecordRead,
		Update: resourceGlesysDNSDomainRecordUpdate,
		Delete: resourceGlesysDNSDomainRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host": {
				Type:     schema.TypeString,
				Required: true,
			},

			"recordid": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceGlesysDNSDomainRecordCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.AddRecordParams{
		Data:       d.Get("data").(string),
		DomainName: d.Get("domain").(string),
		Host:       d.Get("host").(string),
		Type:       d.Get("type").(string),
		TTL:        d.Get("ttl").(int),
	}

	record, err := client.DNSDomains.AddRecord(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error adding record \"%s\": %v", params.Data, err)
	}

	// Set the Id to domain.ID
	id := strconv.Itoa(record.RecordID)
	d.SetId(id)

	return resourceGlesysDNSDomainRecordRead(d, m)
}

func resourceGlesysDNSDomainRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	domain := d.Get("domain").(string)
	myId := d.Get("recordid").(int)
	records, err := client.DNSDomains.ListRecords(context.Background(), domain)

	if err != nil {
		fmt.Errorf("Domain not found: %v\n", err)
		d.SetId("")
		return nil
	}

	for _, record := range *records {
		if record.RecordID == myId {
			d.Set("domain", record.DomainName)
			d.Set("data", record.Data)
			d.Set("host", record.Host)
			d.Set("ttl", record.TTL)
			d.Set("type", record.Type)
		}
	}

	return nil
}

func resourceGlesysDNSDomainRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	myID := d.Id()
	recordID, errid := strconv.Atoi(myID)
	if errid != nil {
		return fmt.Errorf("Id must be converted to integer: %v", errid)
	}
	params := glesys.UpdateRecordParams{RecordID: recordID}

	if d.HasChange("data") {
		params.Data = d.Get("data").(string)
	}

	if d.HasChange("host") {
		params.Host = d.Get("host").(string)
	}

	if d.HasChange("ttl") {
		params.TTL = d.Get("ttl").(int)
	}

	if d.HasChange("type") {
		params.Type = d.Get("type").(string)
	}

	_, err := client.DNSDomains.UpdateRecord(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error updating record: %v", err)
	}

	return resourceGlesysDNSDomainRecordRead(d, m)
}

func resourceGlesysDNSDomainRecordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	recordID, errid := strconv.Atoi(d.Id())
	if errid != nil {
		return fmt.Errorf("Id must be converted to integer: %v", errid)
	}

	err := client.DNSDomains.DeleteRecord(context.Background(), recordID)
	if err != nil {
		return fmt.Errorf("Error deleting domain record: %v", err)
	}
	d.SetId("")
	return nil
}

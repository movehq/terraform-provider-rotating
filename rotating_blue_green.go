// resource_server.go
package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func blueGreen() *schema.Resource {
	return &schema.Resource{
		Create: blueGreenCreate,
		Read:   blueGreenRead,
		Update: blueGreenUpdate,
		Delete: blueGreenDelete,

		Schema: map[string]*schema.Schema{
			"rotate_after_minutes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rotate_after_hours": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rotate_after_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotate_timestamp": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"blue_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"green_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getTimeIncrease(d *schema.ResourceData) (time.Duration, error) {
	rotateMinutes, minExists := d.GetOk("rotate_after_minutes")
	rotateHours, hoursExists := d.GetOk("rotate_after_hours")
	rotateDays, daysExists := d.GetOk("rotate_after_days")

	if !minExists && !hoursExists && !daysExists {
		return 0, errors.New("one of rotate_after_minutes, rotate_after_hours, or rotate_after_days must be set")
	}

	if (minExists && hoursExists) || (minExists && daysExists) || (hoursExists && daysExists) || (minExists && hoursExists && daysExists) {
		return 0, errors.New("only one of rotate_after_minutes, rotate_after_hours, or rotate_after_days cant be set")
	}

	timeToAdd := time.Minute * 1
	if minExists {
		timeToAdd = time.Minute * time.Duration(rotateMinutes.(int))
	} else if hoursExists {
		timeToAdd = time.Hour * time.Duration(rotateHours.(int))
	} else if daysExists {
		timeToAdd = (time.Hour * 24) * time.Duration(rotateDays.(int))
	}

	return timeToAdd, nil
}

func blueGreenCreate(d *schema.ResourceData, m interface{}) error {
	timeToAdd, err := getTimeIncrease(d)
	if err != nil {
		return err
	}

	d.SetId(uuid.New().String())
	err = d.Set("blue_uuid", uuid.New().String())
	if err != nil {
		return err
	}
	err = d.Set("green_uuid", uuid.New().String())
	if err != nil {
		return err
	}
	err = d.Set("active", "blue")
	if err != nil {
		return err
	}

	err = d.Set("rotate_timestamp", int(time.Now().Add(timeToAdd).Unix()))
	if err != nil {
		return err
	}

	return blueGreenRead(d, m)
}

func blueGreenRead(d *schema.ResourceData, m interface{}) error {
	timeToAdd, err := getTimeIncrease(d)
	if err != nil {
		return err
	}
	timestamp := d.Get("rotate_timestamp").(int)
	active := d.Get("active")

	if time.Now().Unix() > int64(timestamp) {
		if active == "green" {
			err := d.Set("active", "blue")
			if err != nil {
				return err
			}
			err = d.Set("blue_uuid", uuid.New().String())
			if err != nil {
				return err
			}
		} else {
			err := d.Set("active", "green")
			if err != nil {
				return err
			}
			err = d.Set("green_uuid", uuid.New().String())
			if err != nil {
				return err
			}
		}
		err := d.Set("rotate_timestamp", int(time.Now().Add(timeToAdd).Unix()))
		if err != nil {
			return err
		}
	}

	return nil
}

func blueGreenUpdate(d *schema.ResourceData, m interface{}) error {
	_, err := getTimeIncrease(d)
	if err != nil {
		return err
	}
	return blueGreenRead(d, m)
}

func blueGreenDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

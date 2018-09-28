package models

import (
	"time"

	"github.com/ilibs/gosql"
)

type Options struct {
	Id          int    `form:"id" json:"id" db:"id"`
	OptionKey   string `form:"option_key" json:"option_key" db:"option_key"`
	OptionValue string `form:"option_value" json:"option_value" db:"option_value"`
}

func (o *Options) DbName() string {
	return "default"
}

func (o *Options) TableName() string {
	return "options"
}

func (o *Options) PK() string {
	return "id"
}

func (o *Options) AfterChange() {
	Cache.Delete("options")
}

func GetOptions() (map[string]string, error) {
	if v, ok := Cache.Get("options"); ok {
		return v.(map[string]string), nil
	}

	var options = make([]*Options, 0)
	err := gosql.Model(&options).All()
	if err != nil {
		return nil, err
	}

	options2 := make(map[string]string)
	for _, v := range options {
		options2[v.OptionKey] = v.OptionValue
	}
	Cache.Set("options", options2, 1*time.Hour)

	return options2, nil
}

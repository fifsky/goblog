package models

type Options struct {
	Id    uint `xorm:"pk"`
	Key   string `xorm:"varchar(100) notnull unique"`
	Value string `xorm:"varchar(200) notnull"`
}

func (o *Options) GetOptions() (map[string]string, error) {
	var options = make([]Options, 0)
	err := orm.Find(&options)
	options2 := make(map[string]string)

	if err != nil {
		return nil, err
	}

	for _, v := range options {
		options2[v.Key] = v.Value
	}

	return options2, nil
}

func (o *Options) Insert() (int64, error) {
	affected, err := orm.Insert(o)
	return affected, err
}

func (o *Options) Update() (int64, error) {
	affected, err := orm.Id(o.Id).Update(o)
	return affected, err
}

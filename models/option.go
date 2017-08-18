package models

type Options struct {
	Id    uint `xorm:"pk"`
	Key   string `xorm:"varchar(100) notnull unique"`
	Value string `xorm:"varchar(200) notnull"`
}

func (this * Options) GetOptions() (map[string]string, error) {
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

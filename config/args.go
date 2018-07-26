package config

import (
	"os"
)

var ExtArgs map[string]string
var Args []string

func argsInit() {
	ExtArgs = make(map[string]string)
	Args = ParseArgs("config", "debug", "show-sql")
}

func ParseArgs(excludes ...string) []string {
	rs := make([]string, 0)
	for _, arg := range os.Args[1:] {
		is_find := false
		for _, ext := range excludes {
			prefix := "--" + ext + "="
			len_prefix := len(prefix)
			if len(arg) > len_prefix && prefix == arg[0:len_prefix] {
				is_find = true
				ExtArgs[ext] = arg[len_prefix:]
			}
		}

		if !is_find {
			rs = append(rs, arg)
		}
	}
	return rs
}

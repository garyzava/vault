package main

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	ess := make(map[string]bool)

	tests := []struct {
		root         string
		in_pattern   string
		out_pathlets []pathlet
	}{
		{"sys", "rekey/backup", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "rekey/backup$", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "auth/(?P<path>.+?)/tune$", []pathlet{{"/sys/auth/{path}/tune", setMaker("path")}}},
		{"sys", "auth/(?P<path>.+?)/tune/(?P<more>.*?)$", []pathlet{{"/sys/auth/{path}/tune/{more}", setMaker("path", "more")}}},
		{"sys", "tools/hash(/(?P<urlalgorithm>.+))?", []pathlet{
			{"/sys/tools/hash/{urlalgorithm}", setMaker("urlalgorithm")},
			{"/sys/tools/hash", setMaker()},
		}},

		/* optional elements
		tools/hash(/(?P<urlalgorithm>.+))?
		{"sys", "leases/lookup/(?P<prefix>.+?)?", []pathlet{
			{"/sys/leases/lookup", []string{"path", "more"}},
			{"/sys/leases/lookup/{prefix}", []string{"prefix"}},
		}},
		*/
		//"leases/lookup/(?P<prefix>.+?)?"
		//"(leases/)?renew"
	}
	/*
		root := "sys"
		pat := "rekey/backup"
		exp := "/sys/rekey/backup"
		out := parsePattern(root, pat)[0].pattern
	*/
	for i, test := range tests {
		out := parsePattern(test.root, test.in_pattern)
		if !reflect.DeepEqual(out, test.out_pathlets) {
			t.Fatalf("Test %d: Expected %v got %v", i, test.out_pathlets, out)
		}
	}
}

func setMaker(strings ...string) map[string]bool {
	ret := make(map[string]bool)

	for _, s := range strings {
		ret[s] = true
	}

	return ret
}

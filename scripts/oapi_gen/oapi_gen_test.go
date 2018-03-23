package main

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	ess := []string{}

	tests := []struct {
		root         string
		in_pattern   string
		out_pathlets []pathlet
	}{
		{"sys", "rekey/backup", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "rekey/backup$", []pathlet{{"/sys/rekey/backup", ess}}},
		{"sys", "auth/(?P<path>.+?)/tune$", []pathlet{{"/sys/auth/{path}/tune", []string{"path"}}}},
		{"sys", "auth/(?P<path>.+?)/tune/(?P<more>.*?)$", []pathlet{{"/sys/auth/{path}/tune/{more}", []string{"path", "more"}}}},
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

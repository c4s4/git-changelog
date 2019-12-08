package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseGitDate(t *testing.T) {
	date, err := time.Parse(dateFormat, "Tue Dec 3 04:37:12 2019 +0100")
	if err != nil {
		t.Errorf("There was an error parsing date: %s", err)
	}
	iso := date.Format(time.RFC3339)
	expected := "2019-12-03T04:37:12+01:00"
	if iso != expected {
		t.Errorf("Bad date parsing: %s but expected %s", iso, expected)
	}
}

func TestParseGitLogs(t *testing.T) {
	logs := `
commit 44be5d95e4a919f229b7867a464437eb259396e3 (HEAD, tag: 1.3.7, origin/master, origin/HEAD, master)
Author: Michel Casabianca <casa@sweetohm.net>
Date:   Thu Dec 5 21:34:21 2019 +0100

	Added completion on templates and themes

commit f7578da9c4f05617cd167d31b6aefc52d240f461
Author: Michel Casabianca <casa@sweetohm.net>
Date:   Tue Dec 3 10:50:34 2019 +0100

	Fixed unit test

commit f2922957684f9c27bf710e99be374c7394843990 (tag: 1.3.2)
Merge: 0dcfb87 c871546
Author: Michel Casabianca <casa@sweetohm.net>
Date:   Sun May 12 20:01:19 2019 +0200

	Release 1.3.2

`
	expected := []commit{
		{
			ID:      "44be5d95e4a919f229b7867a464437eb259396e3",
			Tags:    []string{"1.3.7"},
			Date:    "2019-12-05",
			Message: "Added completion on templates and themes",
		},
		{
			ID:      "f7578da9c4f05617cd167d31b6aefc52d240f461",
			Tags:    nil,
			Date:    "2019-12-03",
			Message: "Fixed unit test",
		},
		{
			ID:      "f2922957684f9c27bf710e99be374c7394843990",
			Tags:    []string{"1.3.2"},
			Date:    "2019-05-12",
			Message: "Release 1.3.2",
		},
	}
	actual, err := parseGitLogs([]byte(logs))
	if err != nil {
		t.Errorf("Error parsing logs: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Bad result: %#v", actual)
	}
}

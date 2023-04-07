package field

import "testing"

func TestSetSort(t *testing.T) {
	var fs FieldSet
	fs = append(fs, &Field{Key: "a"})
	fs = append(fs, &Field{Key: "e"})
	fs = append(fs, &Field{Key: "c"})
	fs.Sort()

	var s string
	for _, v := range fs {
		s += v.Key
	}
	if s != "ace" {
		t.Fatal()
	}

}

func TestSetRemove(t *testing.T) {
	var fs FieldSet
	fs = append(fs, &Field{Key: "a"})
	fs = append(fs, &Field{Key: "e"})
	fs = append(fs, &Field{Key: "c"})

	if fs.Remove("e") == nil {
		t.Fatal()
	}
	if fs.Remove("e") != nil {
		t.Fatal()
	}
	var s string
	for _, v := range fs {
		s += v.Key
	}
	if s != "ac" {
		t.Fatal()
	}
}

package unixtime

import (
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	ut := Now()
	nt := time.Now()
	_, off := nt.Zone()
	if ut.Seconds() != nt.Unix()-int64(off) {
		t.Fatal(ut.Seconds() - nt.Unix())
	}
}

func TestMillis(t *testing.T) {
	ut := Now()
	nt := time.Now()
	_, off := nt.Zone()
	if ut.Millis() != (nt.UnixNano()/1e6 - int64(off)*1e3) {
		t.Fatal(ut.Millis()-nt.UnixNano()/1e6, off)
	}
}

func TestZone(t *testing.T) {

	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	ct, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-7-12 14:45:33", CST)
	at, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-7-12 02:45:33", zone)     //
	gt, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-7-12 06:45:33", time.UTC) // GMT/UTC

	if From(ct) != From(at) {
		t.Fatal(ct, at)
	}
	if From(ct) != From(gt) {
		t.Fatal(ct, at)
	}
}

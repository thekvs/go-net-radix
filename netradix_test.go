package netradix

import (
	"testing"
)

type TestResult struct {
	ip     string
	udata  string
	result bool
}

type TestData struct {
	network string
	udata   string
}

func TestSearchExact(t *testing.T) {
	rtree, err := NewNetRadixTree()
	if err != nil {
		t.Errorf("couldn't create structure")
	}
	defer rtree.Close()

	initial := []TestData{
		{"217.72.192.0/20", "UDATA1"},
		{"217.72.195.0/24", "UDATA2"},
		{"195.161.113.74/32", "UDATA3"},
		{"172.16.2.2", "UDATA4"},
		{"10.42.0.0/16", "UDATA5"},
		{"2001:220::/35", "UDATA6"}}

	for _, value := range initial {
		if err := rtree.Add(value.network, value.udata); err != nil {
			t.Errorf("internal error %v", err)
		}
	}

	expected := []TestResult{
		{"217.72.192.0/20", "UDATA1", true},
		{"217.72.195.42", "", false},
		{"195.161.113.74", "UDATA3", true},
		{"172.16.2.2", "UDATA4", true},
		{"15.161.13.75", "", false},
		{"10.42.1.0/24", "", false},
		{"10.42.1.8", "", false},
		{"2001:220::/128", "", false},
		{"2001:220::/35", "UDATA6", true}}

	for _, value := range expected {
		status, udata, err := rtree.SearchExact(value.ip)
		if err != nil {
			t.Errorf("internal error %v", err)
		}
		if status != value.result {
			t.Errorf("unexpected result for ip %v", value.ip)
		}
		if status {
			if udata != value.udata {
				t.Errorf("unexpected result for ip %v: %v", value.ip, udata)
			}
		}
	}
}

func TestSearchBest(t *testing.T) {
	rtree, err := NewNetRadixTree()
	if err != nil {
		t.Errorf("couldn't create structure")
	}
	defer rtree.Close()

	initial := []TestData{
		{"217.72.192.0/20", "UDATA1"},
		{"217.72.195.0/24", "UDATA2"},
		{"195.161.113.74/32", "UDATA3"},
		{"172.16.2.2", "UDATA4"},
		{"10.42.0.0/16", "UDATA5"},
		{"2001:220::/35", "UDATA6"}}

	expected := []TestResult{
		{"217.72.192.1", "UDATA1", true},
		{"217.72.195.42", "UDATA2", true},
		{"195.161.113.74", "UDATA3", true},
		{"172.16.2.2", "UDATA4", true},
		{"15.161.13.75", "", false},
		{"10.42.1.0/24", "UDATA5", true},
		{"10.42.1.8", "UDATA5", true},
		{"2001:220::/128", "UDATA6", true}}

	for _, value := range initial {
		if err := rtree.Add(value.network, value.udata); err != nil {
			t.Errorf("internal error %v", err)
		}
	}

	for _, value := range expected {
		status, udata, err := rtree.SearchBest(value.ip)
		if err != nil {
			t.Errorf("internal error %v", err)
		}
		if status != value.result {
			t.Errorf("unexpected result for ip %v", value.ip)
		}
		if status {
			if udata != value.udata {
				t.Errorf("unexpected result for ip %v: %v", value.ip, udata)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	rtree, err := NewNetRadixTree()
	if err != nil {
		t.Errorf("couldn't create structure")
	}
	defer rtree.Close()

	initial := []TestData{
		{"217.72.192.0/20", "UDATA1"},
		{"217.72.195.0/24", "UDATA2"},
		{"195.161.113.74/32", "UDATA3"},
		{"172.16.2.2", "UDATA4"},
		{"10.42.0.0/16", "UDATA5"}}

	expected := []TestResult{
		{"217.72.192.0/20", "UDATA1", true},
		{"195.161.113.74", "UDATA3", true},
		{"172.16.2.2", "UDATA4", true}}

	for _, value := range initial {
		if err := rtree.Add(value.network, value.udata); err != nil {
			t.Errorf("internal error %v", err)
		}
	}

	for _, value := range expected {
		status, udata, err := rtree.SearchExact(value.ip)
		if err != nil {
			t.Errorf("internal error %v", err)
		}
		if status != value.result {
			t.Errorf("unexpected result for ip %v", value.ip)
		}
		if status {
			if udata != value.udata {
				t.Errorf("unexpected result for ip %v: %v", value.ip, udata)
			}
		}
	}

	for _, value := range initial {
		rtree.Remove(value.network)
	}

	for _, value := range expected {
		status, _, err := rtree.SearchExact(value.ip)
		if err != nil {
			t.Errorf("internal error %v", err)
		}
		if status != (!value.result) {
			t.Errorf("unexpected result for ip %v", value.ip)
		}
	}
}

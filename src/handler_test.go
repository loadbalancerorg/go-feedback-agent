package main

import "testing"

func TestGetSystemUtilization(t *testing.T) {
	results := getSystemUtilization(80, 0.5, 80, 0)
	if results != 40 {
		t.Errorf("Incorect cpu calculation")
	}

	results = getSystemUtilization(80, 0, 80, 0.5)
	if results != 40 {
		t.Errorf("Incorect ram calculation")
	}

	results = getSystemUtilization(80, 0.5, 80, 0.5)
	if results != 40 {
		t.Errorf("Incorect cpu/ram calculation")
	}

	results = getSystemUtilization(100, 1, 80, 1)
	if results != 100 {
		t.Errorf("Incorect cpu calculation full load")
	}

	results = getSystemUtilization(80, 1, 100, 1)
	if results != 100 {
		t.Errorf("Incorect ram calculation full load")
	}
}

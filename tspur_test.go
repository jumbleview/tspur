package main

import (
	"os"
	"testing"
)

const storage = "demo/cribsheet"
const storageBis = "demo/cribsheetBis"
const pswd = "password"

func TestDecrypt(t *testing.T) {
	_, err := DecryptFile(storage, pswd)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestSave(t *testing.T) {
	data, err := DecryptFile(storage, pswd)
	if err != nil {
		t.Errorf("Decrypt inital file: %s", err)
	}
	var tspr Spur
	tspr.AttachData(data, pswd)
	tspr.cribName = storageBis
	tspr.Save()
	dataBis, err := DecryptFile(storageBis, pswd)
	if err != nil {
		t.Errorf("Decrypt bis file: %s", err)
	}
	if len(data) != len(dataBis) {
		t.Errorf("Different sizes of storages %d!=%d\n", len(data), len(dataBis))
	}

	s := string(data)
	s1 := string(dataBis)
	if s != s1 {
		t.Errorf("Different data %s\n%s\n", s, s1)
	}
	// Cleanup bis file
	err = os.Remove(storageBis)
	if err != nil {
		t.Errorf("Removing bis file: %s", err)
	}
}

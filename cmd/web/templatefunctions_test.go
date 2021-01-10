/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import "testing"

func TestSub(t *testing.T) {
	res := Sub(3, 2)
	if res != 1 {
		t.Error("This is dumb")
	}
}

func TestAdd(t *testing.T) {
	res := Add(1, 2)
	if res != 3 {
		t.Error("Do i really have to do this?")
	}
}

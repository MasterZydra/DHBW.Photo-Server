/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import "testing"

func TestUsernameValidation(t *testing.T) {
	allowedUnames := [6]string{"TestUser", "User123", "max", "ana", "manuela", "robert"}
	forbiddenUnames := [6]string{"%TestUser%", "_User123_", "#max", "(ana)", "<manuela>", "!robert!"}

	for _, name := range allowedUnames {
		err := validateUsername(name)

		if err != nil {
			t.Error("Error while validating allowed username '" + name + "'")
		}
	}

	for _, name := range forbiddenUnames {
		err := validateUsername(name)

		if err == nil {
			t.Error("Error while validating forbidden username '" + name + "'")
		}
	}
}

func TestHostValidation(t *testing.T) {
	allowedHosts := [3]string{"localhost:8080", "https://localhost:4443", "http://192.168.2.162:8888/"}
	forbiddenHosts := [3]string{"localhost", "1234", "192.168.2.162"}

	for _, host := range allowedHosts {
		err := validateHost(host)

		if err != nil {
			t.Error("Error while validating allowed host '" + host + "'")
		}
	}

	for _, host := range forbiddenHosts {
		err := validateHost(host)

		if err == nil {
			t.Error("Error while validating forbidden host '" + host + "'")
		}
	}
}

func ExampleWriteMessage() {
	WriteMessage("Test message")
	// Output:
	// --------------
	// Test message
	// --------------
}

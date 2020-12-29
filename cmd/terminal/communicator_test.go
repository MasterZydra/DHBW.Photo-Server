package main

import "testing"

func TestUsernameValidation (t *testing.T) {
	allowedUnames := [6]string{"TestUser", "User123", "Max", "Ana", "manuela", "robert"}
	forbiddenUnames := [6]string{"%TestUser%", "_User123_", "#Max", "(Ana)", "<manuela>", "!robert!"}

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

func TestHostValidation (t *testing.T) {
	allowedHosts := [6]string{"TestUser", "User123", "Max", "Ana", "manuela", "robert"}
	forbiddenHosts := [6]string{"%TestUser%", "_User123_", "#Max", "(Ana)", "<manuela>", "!robert!"}

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
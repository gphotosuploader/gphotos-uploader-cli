package app_test

import "testing"

func TestStart(t *testing.T) {
	// Should success
	// Should fail when configuration doesn't exist.
	// Should fail when file tracker fails to start.
	// Should fail when token manager fails to start.
	// Should fail when uploads session tracker fails to start.
}

func TestApp_Stop(t *testing.T) {
	// Should success
	// Should fail when file tracker fails to close.
	// Should fail when token manager fails to close.
	// Should fail when uploads session tracker fails to close.
}

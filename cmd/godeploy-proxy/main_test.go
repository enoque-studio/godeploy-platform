package main

import "testing"

func TestGetenv(t *testing.T) {
	t.Setenv("GODEPLOY_TEST_VALUE", "  configured  ")
	if got := getenv("GODEPLOY_TEST_VALUE", "fallback"); got != "configured" {
		t.Fatalf("getenv() = %q, want configured", got)
	}

	t.Setenv("GODEPLOY_TEST_VALUE", "   ")
	if got := getenv("GODEPLOY_TEST_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("getenv() with whitespace = %q, want fallback", got)
	}
}

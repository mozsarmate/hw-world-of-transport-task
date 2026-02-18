package cli_test

import (
	"testing"

	"mate/world-of-transport/internal/cli"
)

func TestParseArgs(t *testing.T) {
	// Valid input
	params, err := cli.ParseArgs([]string{"51.5", "-0.1", "50"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params.Lat != 51.5 || params.Lon != -0.1 || params.DistanceKm != 50 {
		t.Errorf("got %+v, want lat=51.5 lon=-0.1 dist=50", params)
	}

	// Wrong number of args
	_, err = cli.ParseArgs([]string{"51.5", "-0.1"})
	if err == nil {
		t.Error("expected error with 2 args, got nil")
	}

	// Invalid latitude
	_, err = cli.ParseArgs([]string{"not-a-number", "0", "50"})
	if err == nil {
		t.Error("expected error with invalid latitude")
	}

	// Out of range
	_, err = cli.ParseArgs([]string{"100", "0", "50"})
	if err == nil {
		t.Error("expected error with lat > 90")
	}
}

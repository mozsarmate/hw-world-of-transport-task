package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"mate/world-of-transport/internal/geo"
)

type Params struct {
	Lat        float64
	Lon        float64
	DistanceKm float64
}

func ParseArgs(args []string) (*Params, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("expected 3 arguments, got %d", len(args))
	}

	lat, err := parseFloatArg(args[0], "latitude", -90, 90)
	if err != nil {
		return nil, err
	}

	lon, err := parseFloatArg(args[1], "longitude", -180, 180)
	if err != nil {
		return nil, err
	}

	dist, err := parseFloatArg(args[2], "distance_km", 0, 20000)
	if err != nil {
		return nil, err
	}

	return &Params{Lat: lat, Lon: lon, DistanceKm: dist}, nil
}

func parseFloatArg(s, name string, min, max float64) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", name)
	}
	if value < min || value > max {
		return 0, fmt.Errorf("%s must be between %.0f and %.0f", name, min, max)
	}
	return value, nil
}

func PrintHubs(hubs []geo.Hub, params *Params) {
	fmt.Printf("Transport hubs found within %.2f km of (%.4f, %.4f)\n\n",
		params.DistanceKm, params.Lat, params.Lon)

	if len(hubs) == 0 {
		fmt.Println("We did not find any hubs.")
		return
	}

	// Tabwriter for aligned columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "#\tName\tLatitude\tLongitude\tDistance (km)")
	fmt.Fprintln(w, "-\t----\t--------\t---------\t-------------")

	for i, hub := range hubs {
		fmt.Fprintf(w, "%d\t%s\t%.6f\t%.6f\t%.2f\n",
			i+1, hub.Name, hub.Lat, hub.Lon, hub.DistanceKm)
	}

	w.Flush()
	fmt.Printf("\n%d hub(s) found.\n", len(hubs))
}

package main

import (
	"fmt"
	"log"
	"os"

	"mate/world-of-transport/internal/cli"
	"mate/world-of-transport/internal/cloudant"
)

func main() {
	params, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		os.Exit(1)
	}

	client, err := cloudant.NewClient()
	if err != nil {
		log.Fatalf("Failed during Cloudant client initialisation: %v", err)
	}

	hubs, err := client.FindHubsWithinDistance(params.Lat, params.Lon, params.DistanceKm)
	if err != nil {
		log.Fatalf("Failed to fetch transport hubs: %v", err)
	}

	cli.PrintHubs(hubs, params)
}

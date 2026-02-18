package cloudant

import (
	"fmt"
	"sort"

	"mate/world-of-transport/internal/geo"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

const (
	cloudantURL  = "https://mikerhodes.cloudant.com"
	databaseName = "airportdb"
	designDoc    = "view1"
	searchIndex  = "geo"

	maxReturned = 200
)

// Wraps the Cloudant SDK
type Client struct {
	service *cloudantv1.CloudantV1
}

func NewClient() (*Client, error) {
	authenticator := &core.NoAuthAuthenticator{}

	service, err := cloudantv1.NewCloudantV1(&cloudantv1.CloudantV1Options{
		URL:           cloudantURL,
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating Cloudant service: %w", err)
	}

	return &Client{service: service}, nil
}

func getQuery(bbox geo.BoundingBox) string {
	return fmt.Sprintf("lat:[%.6f TO %.6f] AND lon:[%.6f TO %.6f]",
		bbox.MinLat, bbox.MaxLat, bbox.MinLon, bbox.MaxLon)
}

func (c *Client) fetchData(query string) ([]cloudantv1.SearchResultRow, error) {
	var allRows []cloudantv1.SearchResultRow
	var bookmark string

	for {
		opts := c.service.NewPostSearchOptions(databaseName, designDoc, searchIndex, query)
		opts.SetLimit(maxReturned)
		opts.SetIncludeDocs(false)

		if bookmark != "" {
			opts.SetBookmark(bookmark)
		}

		result, _, err := c.service.PostSearch(opts)

		if err != nil {
			return nil, fmt.Errorf("querying Cloudant search: %w", err)
		}

		if len(result.Rows) == 0 {
			break
		}

		allRows = append(allRows, result.Rows...)

		if len(result.Rows) < maxReturned {
			break
		}

		if result.Bookmark == nil || *result.Bookmark == bookmark {
			break
		}
		bookmark = *result.Bookmark
	}

	return allRows, nil
}

func (c *Client) FindHubsWithinDistance(lat, lon, distanceKm float64) ([]geo.Hub, error) {
	bbox := geo.GetBoundingBox(lat, lon, distanceKm)
	query := getQuery(bbox)

	rows, err := c.fetchData(query)
	if err != nil {
		return nil, err
	}

	hubs := processAndFilter(rows, lat, lon, distanceKm)

	sort.Slice(hubs, func(h1, h2 int) bool {
		return hubs[h1].DistanceKm < hubs[h2].DistanceKm
	})

	return hubs, nil
}

func processAndFilter(rows []cloudantv1.SearchResultRow, lat, lon, distanceKm float64) []geo.Hub {
	hubs := make([]geo.Hub, 0, len(rows))
	seen := make(map[string]bool)

	for _, row := range rows {
		hubLat, hubLon, name, ok := extractFields(row)

		//deduplication (appears in London for example)
		key := fmt.Sprintf("%s|%.6f,%.6f", name, hubLat, hubLon)
		curDistance := geo.HaversineDistance(lat, lon, hubLat, hubLon)

		if ok && !seen[key] && curDistance <= distanceKm {

			hubs = append(hubs, geo.Hub{
				Name:       name,
				Lat:        hubLat,
				Lon:        hubLon,
				DistanceKm: curDistance,
			})
			seen[key] = true
		}
	}

	return hubs
}

func extractFields(row cloudantv1.SearchResultRow) (lat, lon float64, name string, ok bool) {
	latVal, hasLat := row.Fields["lat"]
	lonVal, hasLon := row.Fields["lon"]
	nameVal, hasName := row.Fields["name"]

	if !hasLat || !hasLon || !hasName {
		return 0, 0, "", false
	}

	lat, ok = latVal.(float64)
	if !ok {
		return 0, 0, "", false
	}
	lon, ok = lonVal.(float64)
	if !ok {
		return 0, 0, "", false
	}
	name, ok = nameVal.(string)
	if !ok {
		return 0, 0, "", false
	}

	return lat, lon, name, true
}

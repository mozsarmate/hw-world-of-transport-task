# World of Transport \[HW\]

A CLI tool that queries a public Cloudant database of transport hubs and lists
those within a given distance of a specified location, sorted by ascending distance.

## Install, Build, Run

To run this project, you will need GO to be installed (1.21 or later)

```bash
# Clone the repository
git clone https://github.com/mozsarmate/hw-world-of-transport-task.git
cd hw-world-of-transport-task

# Download and install the dependencies
go mod tidy

# Build the binary
go build -o world-of-transport .

# Run the program
./world-of-transport <latitude> <longitude> <distance_km>
```

| Argument       | Description                                 |
|----------------|---------------------------------------------|
| `latitude`     | Decimal degrees, −90 to 90                  |
| `longitude`    | Decimal degrees, −180 to 180                |
| `distance_km`  | Search radius in kilometres (must be > 0)   |

### Examples

Find transport hubs within 50 km of central London:

```bash
./world-of-transport 51.5 -0.1 50
```

Find hubs within 100 km of Budapest:

```bash
./world-of-transport 47.5 19.0 100
```

Find hubs within 0 km of Kiruna (Sweden):

```bash
./world-of-transport 67.852414 20.242500 0
```

## Testing
I have also added some basic tests to verify the CLI's argument parsing. These can be run with:
```bash
go test ./... 
```
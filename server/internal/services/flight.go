package services

import (
	"context"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/repository"
	"github.com/skyquest/server/internal/websocket"
	"github.com/skyquest/server/pkg/aviation"
)

type FlightService struct {
	client      *aviation.Client
	redis       *repository.RedisClient
	flights     []models.Flight
	flightsMux  sync.RWMutex
	airports    map[string]models.Airport
	airportsMux sync.RWMutex
}

func NewFlightService(client *aviation.Client, redis *repository.RedisClient) *FlightService {
	fs := &FlightService{
		client:   client,
		redis:    redis,
		flights:  make([]models.Flight, 0),
		airports: make(map[string]models.Airport),
	}
	fs.initializeAirports()
	// Load initial flight data immediately (don't wait for polling)
	fs.loadInitialFlights()
	return fs
}

// loadInitialFlights fetches flight data on startup
func (s *FlightService) loadInitialFlights() {
	ctx := context.Background()

	// Try Redis cache first (but don't fail if Redis is down)
	if s.redis != nil {
		cached, err := s.redis.GetCachedFlights(ctx)
		if err == nil && cached != nil && len(cached) > 0 {
			s.updateFlights(cached)
			log.Printf("Loaded %d flights from cache", len(cached))
			return
		}
	}

	// Fetch from API
	flights, err := s.client.GetFlights()
	if err != nil {
		log.Printf("Error fetching initial flights: %v", err)
		return
	}

	// Enrich flights with airport data and coordinates
	s.enrichFlights(flights)

	// Try to cache the results (don't fail if Redis is down)
	if s.redis != nil {
		if err := s.redis.CacheFlights(ctx, flights); err != nil {
			log.Printf("Warning: Could not cache flights: %v", err)
		}
	}

	s.updateFlights(flights)
	log.Printf("Loaded %d initial flights", len(flights))
}

// initializeAirports loads the airport database
func (s *FlightService) initializeAirports() {
	// Major world airports for the game
	airports := []models.Airport{
		{IATA: "JFK", ICAO: "KJFK", Name: "John F. Kennedy International", City: "New York", Country: "USA", Latitude: 40.6413, Longitude: -73.7781},
		{IATA: "LGA", ICAO: "KLGA", Name: "LaGuardia Airport", City: "New York", Country: "USA", Latitude: 40.7769, Longitude: -73.8740},
		{IATA: "EWR", ICAO: "KEWR", Name: "Newark Liberty International", City: "Newark", Country: "USA", Latitude: 40.6895, Longitude: -74.1745},
		{IATA: "LAX", ICAO: "KLAX", Name: "Los Angeles International", City: "Los Angeles", Country: "USA", Latitude: 33.9425, Longitude: -118.4081},
		{IATA: "SFO", ICAO: "KSFO", Name: "San Francisco International", City: "San Francisco", Country: "USA", Latitude: 37.6213, Longitude: -122.3790},
		{IATA: "ORD", ICAO: "KORD", Name: "O'Hare International", City: "Chicago", Country: "USA", Latitude: 41.9742, Longitude: -87.9073},
		{IATA: "MIA", ICAO: "KMIA", Name: "Miami International", City: "Miami", Country: "USA", Latitude: 25.7959, Longitude: -80.2870},
		{IATA: "BOS", ICAO: "KBOS", Name: "Boston Logan International", City: "Boston", Country: "USA", Latitude: 42.3656, Longitude: -71.0096},
		{IATA: "ATL", ICAO: "KATL", Name: "Hartsfield-Jackson Atlanta", City: "Atlanta", Country: "USA", Latitude: 33.6407, Longitude: -84.4277},
		{IATA: "DFW", ICAO: "KDFW", Name: "Dallas/Fort Worth International", City: "Dallas", Country: "USA", Latitude: 32.8998, Longitude: -97.0403},
		{IATA: "SEA", ICAO: "KSEA", Name: "Seattle-Tacoma International", City: "Seattle", Country: "USA", Latitude: 47.4502, Longitude: -122.3088},
		{IATA: "YYZ", ICAO: "CYYZ", Name: "Toronto Pearson International", City: "Toronto", Country: "Canada", Latitude: 43.6777, Longitude: -79.6248},
		{IATA: "YVR", ICAO: "CYVR", Name: "Vancouver International", City: "Vancouver", Country: "Canada", Latitude: 49.1967, Longitude: -123.1815},
		{IATA: "MEX", ICAO: "MMMX", Name: "Mexico City International", City: "Mexico City", Country: "Mexico", Latitude: 19.4361, Longitude: -99.0719},
		{IATA: "LHR", ICAO: "EGLL", Name: "London Heathrow", City: "London", Country: "UK", Latitude: 51.4700, Longitude: -0.4543},
		{IATA: "LGW", ICAO: "EGKK", Name: "London Gatwick", City: "London", Country: "UK", Latitude: 51.1537, Longitude: -0.1821},
		{IATA: "CDG", ICAO: "LFPG", Name: "Charles de Gaulle", City: "Paris", Country: "France", Latitude: 49.0097, Longitude: 2.5479},
		{IATA: "ORY", ICAO: "LFPO", Name: "Paris Orly", City: "Paris", Country: "France", Latitude: 48.7233, Longitude: 2.3795},
		{IATA: "FRA", ICAO: "EDDF", Name: "Frankfurt Airport", City: "Frankfurt", Country: "Germany", Latitude: 50.0379, Longitude: 8.5622},
		{IATA: "MUC", ICAO: "EDDM", Name: "Munich Airport", City: "Munich", Country: "Germany", Latitude: 48.3537, Longitude: 11.7750},
		{IATA: "AMS", ICAO: "EHAM", Name: "Amsterdam Schiphol", City: "Amsterdam", Country: "Netherlands", Latitude: 52.3105, Longitude: 4.7683},
		{IATA: "MAD", ICAO: "LEMD", Name: "Madrid Barajas", City: "Madrid", Country: "Spain", Latitude: 40.4983, Longitude: -3.5676},
		{IATA: "BCN", ICAO: "LEBL", Name: "Barcelona El Prat", City: "Barcelona", Country: "Spain", Latitude: 41.2974, Longitude: 2.0833},
		{IATA: "FCO", ICAO: "LIRF", Name: "Rome Fiumicino", City: "Rome", Country: "Italy", Latitude: 41.8003, Longitude: 12.2389},
		{IATA: "ZRH", ICAO: "LSZH", Name: "Zurich Airport", City: "Zurich", Country: "Switzerland", Latitude: 47.4647, Longitude: 8.5492},
		{IATA: "VIE", ICAO: "LOWW", Name: "Vienna International", City: "Vienna", Country: "Austria", Latitude: 48.1103, Longitude: 16.5697},
		{IATA: "CPH", ICAO: "EKCH", Name: "Copenhagen Airport", City: "Copenhagen", Country: "Denmark", Latitude: 55.6180, Longitude: 12.6560},
		{IATA: "DUB", ICAO: "EIDW", Name: "Dublin Airport", City: "Dublin", Country: "Ireland", Latitude: 53.4264, Longitude: -6.2499},
		{IATA: "IST", ICAO: "LTFM", Name: "Istanbul Airport", City: "Istanbul", Country: "Turkey", Latitude: 41.2753, Longitude: 28.7519},
		{IATA: "DXB", ICAO: "OMDB", Name: "Dubai International", City: "Dubai", Country: "UAE", Latitude: 25.2532, Longitude: 55.3657},
		{IATA: "HKG", ICAO: "VHHH", Name: "Hong Kong International", City: "Hong Kong", Country: "Hong Kong", Latitude: 22.3080, Longitude: 113.9185},
		{IATA: "SIN", ICAO: "WSSS", Name: "Singapore Changi", City: "Singapore", Country: "Singapore", Latitude: 1.3644, Longitude: 103.9915},
		{IATA: "NRT", ICAO: "RJAA", Name: "Narita International", City: "Tokyo", Country: "Japan", Latitude: 35.7720, Longitude: 140.3929},
		{IATA: "HND", ICAO: "RJTT", Name: "Tokyo Haneda", City: "Tokyo", Country: "Japan", Latitude: 35.5494, Longitude: 139.7798},
		{IATA: "ICN", ICAO: "RKSI", Name: "Incheon International", City: "Seoul", Country: "South Korea", Latitude: 37.4691, Longitude: 126.4505},
		{IATA: "PEK", ICAO: "ZBAA", Name: "Beijing Capital International", City: "Beijing", Country: "China", Latitude: 40.0799, Longitude: 116.6031},
		{IATA: "PVG", ICAO: "ZSPD", Name: "Shanghai Pudong International", City: "Shanghai", Country: "China", Latitude: 31.1443, Longitude: 121.8083},
		{IATA: "BKK", ICAO: "VTBS", Name: "Suvarnabhumi Airport", City: "Bangkok", Country: "Thailand", Latitude: 13.6900, Longitude: 100.7501},
		{IATA: "KUL", ICAO: "WMKK", Name: "Kuala Lumpur International", City: "Kuala Lumpur", Country: "Malaysia", Latitude: 2.7456, Longitude: 101.7099},
		{IATA: "DEL", ICAO: "VIDP", Name: "Indira Gandhi International", City: "New Delhi", Country: "India", Latitude: 28.5562, Longitude: 77.1000},
		{IATA: "BOM", ICAO: "VABB", Name: "Chhatrapati Shivaji International", City: "Mumbai", Country: "India", Latitude: 19.0896, Longitude: 72.8656},
		{IATA: "SYD", ICAO: "YSSY", Name: "Sydney Kingsford Smith", City: "Sydney", Country: "Australia", Latitude: -33.9399, Longitude: 151.1753},
		{IATA: "MEL", ICAO: "YMML", Name: "Melbourne Airport", City: "Melbourne", Country: "Australia", Latitude: -37.6690, Longitude: 144.8410},
		{IATA: "AKL", ICAO: "NZAA", Name: "Auckland Airport", City: "Auckland", Country: "New Zealand", Latitude: -37.0082, Longitude: 174.7850},
		{IATA: "DOH", ICAO: "OTHH", Name: "Hamad International", City: "Doha", Country: "Qatar", Latitude: 25.2731, Longitude: 51.6081},
		{IATA: "GRU", ICAO: "SBGR", Name: "São Paulo–Guarulhos International", City: "São Paulo", Country: "Brazil", Latitude: -23.4356, Longitude: -46.4731},
		{IATA: "EZE", ICAO: "SAEZ", Name: "Ministro Pistarini International", City: "Buenos Aires", Country: "Argentina", Latitude: -34.8222, Longitude: -58.5358},
		{IATA: "JNB", ICAO: "FAOR", Name: "O.R. Tambo International", City: "Johannesburg", Country: "South Africa", Latitude: -26.1367, Longitude: 28.2411},
		{IATA: "CAI", ICAO: "HECA", Name: "Cairo International", City: "Cairo", Country: "Egypt", Latitude: 30.1219, Longitude: 31.4056},
	}

	s.airportsMux.Lock()
	defer s.airportsMux.Unlock()
	for _, airport := range airports {
		s.airports[airport.IATA] = airport
	}
}

// StartPolling starts the background flight data polling
func (s *FlightService) StartPolling(hub *websocket.Hub, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial fetch
	s.fetchAndBroadcast(hub)

	for range ticker.C {
		s.fetchAndBroadcast(hub)
	}
}

func (s *FlightService) fetchAndBroadcast(hub *websocket.Hub) {
	ctx := context.Background()

	// Check cache first
	cached, err := s.redis.GetCachedFlights(ctx)
	if err == nil && cached != nil {
		s.updateFlights(cached)
		hub.BroadcastFlights(cached)
		return
	}

	// Fetch from API
	flights, err := s.client.GetFlights()
	if err != nil {
		log.Printf("Error fetching flights: %v", err)
		return
	}

	// Enrich flights with airport data and coordinates
	s.enrichFlights(flights)

	// Cache the results
	if err := s.redis.CacheFlights(ctx, flights); err != nil {
		log.Printf("Error caching flights: %v", err)
	}

	s.updateFlights(flights)
	hub.BroadcastFlights(flights)
}

func (s *FlightService) updateFlights(flights []models.Flight) {
	s.flightsMux.Lock()
	defer s.flightsMux.Unlock()
	s.flights = flights
}

// enrichFlights adds airport data and generates initial coordinates
// The actual position/altitude/speed/heading will be randomized in prepareFlightForDisplay
func (s *FlightService) enrichFlights(flights []models.Flight) {
	for i := range flights {
		// Enrich departure airport data
		if airport, ok := s.GetAirport(flights[i].Departure.IATA); ok {
			flights[i].Departure = airport
		}

		// Enrich arrival airport data
		if airport, ok := s.GetAirport(flights[i].Arrival.IATA); ok {
			flights[i].Arrival = airport
		} else {
			// Airport not in our database - try to extract city from airport name
			if flights[i].Arrival.City == "" && flights[i].Arrival.Name != "" {
				flights[i].Arrival.City = extractCityFromAirportName(flights[i].Arrival.Name)
			}
		}

		// Use DEPARTURE airport for positioning the aircraft
		departureLat := flights[i].Departure.Latitude
		departureLon := flights[i].Departure.Longitude
		arrivalLat := flights[i].Arrival.Latitude
		arrivalLon := flights[i].Arrival.Longitude

		// Fallback if arrival is missing (for heading calculation)
		if arrivalLat == 0 && arrivalLon == 0 {
			arrivalLat = departureLat + 5.0
			arrivalLon = departureLon + 5.0
		}

		// Generate initial random position CENTERED around the DEPARTURE airport
		if departureLat != 0 || departureLon != 0 {
			// Generate small random offset (0.3-1.0 degrees away, ~33-111 km) for close proximity
			distance := 0.3 + rand.Float64()*0.7  // 0.3 to 1.0 degrees (~33-111 km)
			angle := rand.Float64() * 2 * math.Pi // Random angle in radians

			// Calculate aircraft position at random offset from DEPARTURE airport (centered around it)
			flights[i].Latitude = departureLat + distance*math.Cos(angle)
			flights[i].Longitude = departureLon + distance*math.Sin(angle)

			// Calculate heading from aircraft position towards the ARRIVAL airport
			flights[i].Direction = calculateHeading(
				flights[i].Latitude, flights[i].Longitude,
				arrivalLat, arrivalLon,
			)

			// Randomize flight data
			flights[i].Altitude = 28000 + rand.Float64()*10000    // 28,000-38,000 ft
			flights[i].Speed = 420 + rand.Float64()*100           // 420-520 knots
			flights[i].VerticalSpeed = -500 + rand.Float64()*1000 // -500 to +500 ft/min
		}
	}
}

// calculateHeading calculates the bearing/heading from point 1 to point 2 in degrees
func calculateHeading(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Calculate bearing
	x := math.Sin(deltaLon) * math.Cos(lat2Rad)
	y := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(deltaLon)

	bearing := math.Atan2(x, y)

	// Convert to degrees and normalize to 0-360
	bearingDeg := bearing * 180 / math.Pi
	if bearingDeg < 0 {
		bearingDeg += 360
	}

	return bearingDeg
}

// extractCityFromAirportName attempts to extract a city name from airport name
func extractCityFromAirportName(name string) string {
	// Common patterns: "City International", "City Airport", "City-Name Airport"
	suffixes := []string{" International", " Airport", " Intl", " Regional", " Municipal"}
	result := name
	for _, suffix := range suffixes {
		if idx := indexOfStr(result, suffix); idx > 0 {
			result = result[:idx]
		}
	}
	return result
}

func indexOfStr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetFlights returns filtered flights based on difficulty
func (s *FlightService) GetFlights(difficulty models.Difficulty) []models.Flight {
	s.flightsMux.RLock()
	defer s.flightsMux.RUnlock()

	var filtered []models.Flight
	for _, f := range s.flights {
		if s.matchesDifficulty(f, difficulty) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// GetRandomFlights returns n random flights matching difficulty filter
func (s *FlightService) GetRandomFlights(difficulty models.Difficulty, n int) []models.Flight {
	flights := s.GetFlights(difficulty)

	// Fallback: if no flights match difficulty, try without filter
	if len(flights) == 0 {
		log.Printf("No flights match difficulty=%s, getting all flights", difficulty)
		s.flightsMux.RLock()
		flights = make([]models.Flight, len(s.flights))
		copy(flights, s.flights)
		s.flightsMux.RUnlock()
	}

	// Last resort: try to reload flights
	if len(flights) == 0 {
		log.Println("No flights available, attempting to reload...")
		s.loadInitialFlights()
		s.flightsMux.RLock()
		flights = make([]models.Flight, len(s.flights))
		copy(flights, s.flights)
		s.flightsMux.RUnlock()
	}

	if len(flights) == 0 {
		log.Println("Still no flights available after reload attempt")
		return nil
	}

	// Shuffle and take n flights
	rand.Shuffle(len(flights), func(i, j int) {
		flights[i], flights[j] = flights[j], flights[i]
	})

	if n > len(flights) {
		n = len(flights)
	}
	return flights[:n]
}

// GetFlightByID returns a specific flight
func (s *FlightService) GetFlightByID(id string) *models.Flight {
	s.flightsMux.RLock()
	defer s.flightsMux.RUnlock()

	for _, f := range s.flights {
		if f.ID == id {
			return &f
		}
	}
	return nil
}

// GetAirport returns airport info by IATA code
func (s *FlightService) GetAirport(iata string) (models.Airport, bool) {
	s.airportsMux.RLock()
	defer s.airportsMux.RUnlock()
	airport, ok := s.airports[iata]
	return airport, ok
}

// GetAllAirports returns all airports
func (s *FlightService) GetAllAirports() []models.Airport {
	s.airportsMux.RLock()
	defer s.airportsMux.RUnlock()

	airports := make([]models.Airport, 0, len(s.airports))
	for _, a := range s.airports {
		airports = append(airports, a)
	}
	return airports
}

// matchesDifficulty filters flights based on difficulty level
// Easy: Domestic flights only (same country)
// Medium: Short to medium-haul flights (distance < 5000km)
// Hard: All flights including long-haul international
func (s *FlightService) matchesDifficulty(f models.Flight, difficulty models.Difficulty) bool {
	switch difficulty {
	case models.DifficultyEasy:
		// Domestic only - same country
		if f.Departure.Country != f.Arrival.Country {
			return false
		}
	case models.DifficultyMedium:
		// Short to medium-haul flights (distance < 5000km)
		distance := CalculateDistance(
			f.Departure.Latitude, f.Departure.Longitude,
			f.Arrival.Latitude, f.Arrival.Longitude,
		)
		if distance > 5000 {
			return false
		}
	case models.DifficultyHard:
		// All flights allowed
	}

	return true
}

// CalculateDistance calculates the distance between two points using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in km

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

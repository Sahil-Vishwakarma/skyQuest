package aviation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skyquest/server/internal/models"
)

const (
	baseURL = "http://api.aviationstack.com/v1"
)

// Client is the AviationStack API client
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new AviationStack API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AviationStackResponse represents the API response wrapper
type AviationStackResponse struct {
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Data []FlightData `json:"data"`
}

// FlightData represents a single flight from AviationStack
type FlightData struct {
	FlightDate   string `json:"flight_date"`
	FlightStatus string `json:"flight_status"`
	Departure    struct {
		Airport   string  `json:"airport"`
		Timezone  string  `json:"timezone"`
		IATA      string  `json:"iata"`
		ICAO      string  `json:"icao"`
		Terminal  string  `json:"terminal"`
		Gate      string  `json:"gate"`
		Delay     int     `json:"delay"`
		Scheduled string  `json:"scheduled"`
		Estimated string  `json:"estimated"`
		Actual    string  `json:"actual"`
		ActualRwy string  `json:"actual_runway"`
	} `json:"departure"`
	Arrival struct {
		Airport   string  `json:"airport"`
		Timezone  string  `json:"timezone"`
		IATA      string  `json:"iata"`
		ICAO      string  `json:"icao"`
		Terminal  string  `json:"terminal"`
		Gate      string  `json:"gate"`
		Baggage   string  `json:"baggage"`
		Delay     int     `json:"delay"`
		Scheduled string  `json:"scheduled"`
		Estimated string  `json:"estimated"`
		Actual    string  `json:"actual"`
		ActualRwy string  `json:"actual_runway"`
	} `json:"arrival"`
	Airline struct {
		Name string `json:"name"`
		IATA string `json:"iata"`
		ICAO string `json:"icao"`
	} `json:"airline"`
	Flight struct {
		Number     string `json:"number"`
		IATA       string `json:"iata"`
		ICAO       string `json:"icao"`
		Codeshared *struct {
			AirlineName  string `json:"airline_name"`
			AirlineIATA  string `json:"airline_iata"`
			AirlineICAO  string `json:"airline_icao"`
			FlightNumber string `json:"flight_number"`
			FlightIATA   string `json:"flight_iata"`
			FlightICAO   string `json:"flight_icao"`
		} `json:"codeshared"`
	} `json:"flight"`
	Aircraft *struct {
		Registration string `json:"registration"`
		IATA         string `json:"iata"`
		ICAO         string `json:"icao"`
		ICAO24       string `json:"icao24"`
	} `json:"aircraft"`
	Live *struct {
		Updated      string  `json:"updated"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Altitude     float64 `json:"altitude"`
		Direction    float64 `json:"direction"`
		SpeedH       float64 `json:"speed_horizontal"`
		SpeedV       float64 `json:"speed_vertical"`
		IsGround     bool    `json:"is_ground"`
	} `json:"live"`
}

// GetFlights fetches all currently active flights from AviationStack API
func (c *Client) GetFlights() ([]models.Flight, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("AviationStack API key is required - set AVIATIONSTACK_API_KEY environment variable")
	}

	// Fetch active/en-route flights
	url := fmt.Sprintf("%s/flights?access_key=%s&flight_status=active&limit=100", baseURL, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flights: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResponse AviationStackResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.convertToFlights(apiResponse.Data), nil
}

func (c *Client) convertToFlights(raw []FlightData) []models.Flight {
	flights := make([]models.Flight, 0, len(raw))

	for _, r := range raw {
		// Skip flights without proper departure/arrival info
		if r.Departure.IATA == "" || r.Arrival.IATA == "" {
			continue
		}

		// Generate a unique ID for the flight
		flightID := r.Flight.IATA
		if r.Aircraft != nil && r.Aircraft.ICAO24 != "" {
			flightID = r.Aircraft.ICAO24
		}
		if flightID == "" {
			flightID = fmt.Sprintf("%s-%s-%s", r.Flight.IATA, r.Departure.IATA, r.Arrival.IATA)
		}

		flight := models.Flight{
			ID:           flightID,
			ICAO24:       getAircraftICAO24(r.Aircraft),
			Callsign:     r.Flight.ICAO,
			FlightNumber: r.Flight.IATA,
			Status:       r.FlightStatus,
			Departure: models.Airport{
				IATA: r.Departure.IATA,
				ICAO: r.Departure.ICAO,
				Name: r.Departure.Airport,
			},
			Arrival: models.Airport{
				IATA: r.Arrival.IATA,
				ICAO: r.Arrival.ICAO,
				Name: r.Arrival.Airport,
			},
			Aircraft: models.Aircraft{
				IATA:         getAircraftIATA(r.Aircraft),
				ICAO:         getAircraftICAO(r.Aircraft),
				Registration: getAircraftReg(r.Aircraft),
			},
			Airline: models.Airline{
				IATA: r.Airline.IATA,
				ICAO: r.Airline.ICAO,
				Name: r.Airline.Name,
			},
			UpdatedAt: time.Now(),
		}

		// Add live tracking data if available
		if r.Live != nil {
			flight.Latitude = r.Live.Latitude
			flight.Longitude = r.Live.Longitude
			flight.Altitude = r.Live.Altitude
			flight.Direction = r.Live.Direction
			flight.Speed = r.Live.SpeedH
			flight.VerticalSpeed = r.Live.SpeedV
		}

		flights = append(flights, flight)
	}

	return flights
}

// Helper functions to safely access aircraft data
func getAircraftICAO24(a *struct {
	Registration string `json:"registration"`
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	ICAO24       string `json:"icao24"`
}) string {
	if a == nil {
		return ""
	}
	return a.ICAO24
}

func getAircraftIATA(a *struct {
	Registration string `json:"registration"`
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	ICAO24       string `json:"icao24"`
}) string {
	if a == nil {
		return ""
	}
	return a.IATA
}

func getAircraftICAO(a *struct {
	Registration string `json:"registration"`
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	ICAO24       string `json:"icao24"`
}) string {
	if a == nil {
		return ""
	}
	return a.ICAO
}

func getAircraftReg(a *struct {
	Registration string `json:"registration"`
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	ICAO24       string `json:"icao24"`
}) string {
	if a == nil {
		return ""
	}
	return a.Registration
}

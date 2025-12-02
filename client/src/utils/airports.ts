import type { Airport } from '../types'

// Major world airports for autocomplete suggestions
export const airports: Airport[] = [
  { iata: 'JFK', icao: 'KJFK', name: 'John F. Kennedy International', city: 'New York', country: 'USA', latitude: 40.6413, longitude: -73.7781 },
  { iata: 'LGA', icao: 'KLGA', name: 'LaGuardia Airport', city: 'New York', country: 'USA', latitude: 40.7769, longitude: -73.8740 },
  { iata: 'EWR', icao: 'KEWR', name: 'Newark Liberty International', city: 'Newark', country: 'USA', latitude: 40.6895, longitude: -74.1745 },
  { iata: 'LAX', icao: 'KLAX', name: 'Los Angeles International', city: 'Los Angeles', country: 'USA', latitude: 33.9425, longitude: -118.4081 },
  { iata: 'SFO', icao: 'KSFO', name: 'San Francisco International', city: 'San Francisco', country: 'USA', latitude: 37.6213, longitude: -122.3790 },
  { iata: 'ORD', icao: 'KORD', name: "O'Hare International", city: 'Chicago', country: 'USA', latitude: 41.9742, longitude: -87.9073 },
  { iata: 'MIA', icao: 'KMIA', name: 'Miami International', city: 'Miami', country: 'USA', latitude: 25.7959, longitude: -80.2870 },
  { iata: 'BOS', icao: 'KBOS', name: 'Boston Logan International', city: 'Boston', country: 'USA', latitude: 42.3656, longitude: -71.0096 },
  { iata: 'ATL', icao: 'KATL', name: 'Hartsfield-Jackson Atlanta', city: 'Atlanta', country: 'USA', latitude: 33.6407, longitude: -84.4277 },
  { iata: 'DFW', icao: 'KDFW', name: 'Dallas/Fort Worth International', city: 'Dallas', country: 'USA', latitude: 32.8998, longitude: -97.0403 },
  { iata: 'SEA', icao: 'KSEA', name: 'Seattle-Tacoma International', city: 'Seattle', country: 'USA', latitude: 47.4502, longitude: -122.3088 },
  { iata: 'DEN', icao: 'KDEN', name: 'Denver International', city: 'Denver', country: 'USA', latitude: 39.8561, longitude: -104.6737 },
  { iata: 'PHX', icao: 'KPHX', name: 'Phoenix Sky Harbor', city: 'Phoenix', country: 'USA', latitude: 33.4373, longitude: -112.0078 },
  { iata: 'LAS', icao: 'KLAS', name: 'Harry Reid International', city: 'Las Vegas', country: 'USA', latitude: 36.0840, longitude: -115.1537 },
  { iata: 'YYZ', icao: 'CYYZ', name: 'Toronto Pearson International', city: 'Toronto', country: 'Canada', latitude: 43.6777, longitude: -79.6248 },
  { iata: 'YVR', icao: 'CYVR', name: 'Vancouver International', city: 'Vancouver', country: 'Canada', latitude: 49.1967, longitude: -123.1815 },
  { iata: 'YUL', icao: 'CYUL', name: 'Montréal-Trudeau International', city: 'Montreal', country: 'Canada', latitude: 45.4706, longitude: -73.7408 },
  { iata: 'MEX', icao: 'MMMX', name: 'Mexico City International', city: 'Mexico City', country: 'Mexico', latitude: 19.4361, longitude: -99.0719 },
  { iata: 'CUN', icao: 'MMUN', name: 'Cancún International', city: 'Cancún', country: 'Mexico', latitude: 21.0365, longitude: -86.8770 },
  { iata: 'LHR', icao: 'EGLL', name: 'London Heathrow', city: 'London', country: 'UK', latitude: 51.4700, longitude: -0.4543 },
  { iata: 'LGW', icao: 'EGKK', name: 'London Gatwick', city: 'London', country: 'UK', latitude: 51.1537, longitude: -0.1821 },
  { iata: 'STN', icao: 'EGSS', name: 'London Stansted', city: 'London', country: 'UK', latitude: 51.8860, longitude: 0.2389 },
  { iata: 'CDG', icao: 'LFPG', name: 'Charles de Gaulle', city: 'Paris', country: 'France', latitude: 49.0097, longitude: 2.5479 },
  { iata: 'ORY', icao: 'LFPO', name: 'Paris Orly', city: 'Paris', country: 'France', latitude: 48.7233, longitude: 2.3795 },
  { iata: 'FRA', icao: 'EDDF', name: 'Frankfurt Airport', city: 'Frankfurt', country: 'Germany', latitude: 50.0379, longitude: 8.5622 },
  { iata: 'MUC', icao: 'EDDM', name: 'Munich Airport', city: 'Munich', country: 'Germany', latitude: 48.3537, longitude: 11.7750 },
  { iata: 'BER', icao: 'EDDB', name: 'Berlin Brandenburg', city: 'Berlin', country: 'Germany', latitude: 52.3667, longitude: 13.5033 },
  { iata: 'AMS', icao: 'EHAM', name: 'Amsterdam Schiphol', city: 'Amsterdam', country: 'Netherlands', latitude: 52.3105, longitude: 4.7683 },
  { iata: 'MAD', icao: 'LEMD', name: 'Madrid Barajas', city: 'Madrid', country: 'Spain', latitude: 40.4983, longitude: -3.5676 },
  { iata: 'BCN', icao: 'LEBL', name: 'Barcelona El Prat', city: 'Barcelona', country: 'Spain', latitude: 41.2974, longitude: 2.0833 },
  { iata: 'FCO', icao: 'LIRF', name: 'Rome Fiumicino', city: 'Rome', country: 'Italy', latitude: 41.8003, longitude: 12.2389 },
  { iata: 'MXP', icao: 'LIMC', name: 'Milan Malpensa', city: 'Milan', country: 'Italy', latitude: 45.6306, longitude: 8.7281 },
  { iata: 'ZRH', icao: 'LSZH', name: 'Zurich Airport', city: 'Zurich', country: 'Switzerland', latitude: 47.4647, longitude: 8.5492 },
  { iata: 'VIE', icao: 'LOWW', name: 'Vienna International', city: 'Vienna', country: 'Austria', latitude: 48.1103, longitude: 16.5697 },
  { iata: 'CPH', icao: 'EKCH', name: 'Copenhagen Airport', city: 'Copenhagen', country: 'Denmark', latitude: 55.6180, longitude: 12.6560 },
  { iata: 'ARN', icao: 'ESSA', name: 'Stockholm Arlanda', city: 'Stockholm', country: 'Sweden', latitude: 59.6519, longitude: 17.9186 },
  { iata: 'OSL', icao: 'ENGM', name: 'Oslo Gardermoen', city: 'Oslo', country: 'Norway', latitude: 60.1939, longitude: 11.1004 },
  { iata: 'HEL', icao: 'EFHK', name: 'Helsinki-Vantaa', city: 'Helsinki', country: 'Finland', latitude: 60.3172, longitude: 24.9633 },
  { iata: 'DUB', icao: 'EIDW', name: 'Dublin Airport', city: 'Dublin', country: 'Ireland', latitude: 53.4264, longitude: -6.2499 },
  { iata: 'LIS', icao: 'LPPT', name: 'Lisbon Portela', city: 'Lisbon', country: 'Portugal', latitude: 38.7756, longitude: -9.1354 },
  { iata: 'ATH', icao: 'LGAV', name: 'Athens International', city: 'Athens', country: 'Greece', latitude: 37.9364, longitude: 23.9445 },
  { iata: 'IST', icao: 'LTFM', name: 'Istanbul Airport', city: 'Istanbul', country: 'Turkey', latitude: 41.2753, longitude: 28.7519 },
  { iata: 'DXB', icao: 'OMDB', name: 'Dubai International', city: 'Dubai', country: 'UAE', latitude: 25.2532, longitude: 55.3657 },
  { iata: 'AUH', icao: 'OMAA', name: 'Abu Dhabi International', city: 'Abu Dhabi', country: 'UAE', latitude: 24.4330, longitude: 54.6511 },
  { iata: 'DOH', icao: 'OTHH', name: 'Hamad International', city: 'Doha', country: 'Qatar', latitude: 25.2731, longitude: 51.6081 },
  { iata: 'HKG', icao: 'VHHH', name: 'Hong Kong International', city: 'Hong Kong', country: 'Hong Kong', latitude: 22.3080, longitude: 113.9185 },
  { iata: 'SIN', icao: 'WSSS', name: 'Singapore Changi', city: 'Singapore', country: 'Singapore', latitude: 1.3644, longitude: 103.9915 },
  { iata: 'NRT', icao: 'RJAA', name: 'Narita International', city: 'Tokyo', country: 'Japan', latitude: 35.7720, longitude: 140.3929 },
  { iata: 'HND', icao: 'RJTT', name: 'Tokyo Haneda', city: 'Tokyo', country: 'Japan', latitude: 35.5494, longitude: 139.7798 },
  { iata: 'KIX', icao: 'RJBB', name: 'Kansai International', city: 'Osaka', country: 'Japan', latitude: 34.4347, longitude: 135.2441 },
  { iata: 'ICN', icao: 'RKSI', name: 'Incheon International', city: 'Seoul', country: 'South Korea', latitude: 37.4691, longitude: 126.4505 },
  { iata: 'PEK', icao: 'ZBAA', name: 'Beijing Capital International', city: 'Beijing', country: 'China', latitude: 40.0799, longitude: 116.6031 },
  { iata: 'PVG', icao: 'ZSPD', name: 'Shanghai Pudong International', city: 'Shanghai', country: 'China', latitude: 31.1443, longitude: 121.8083 },
  { iata: 'CAN', icao: 'ZGGG', name: 'Guangzhou Baiyun', city: 'Guangzhou', country: 'China', latitude: 23.3924, longitude: 113.2988 },
  { iata: 'TPE', icao: 'RCTP', name: 'Taiwan Taoyuan International', city: 'Taipei', country: 'Taiwan', latitude: 25.0797, longitude: 121.2342 },
  { iata: 'BKK', icao: 'VTBS', name: 'Suvarnabhumi Airport', city: 'Bangkok', country: 'Thailand', latitude: 13.6900, longitude: 100.7501 },
  { iata: 'KUL', icao: 'WMKK', name: 'Kuala Lumpur International', city: 'Kuala Lumpur', country: 'Malaysia', latitude: 2.7456, longitude: 101.7099 },
  { iata: 'CGK', icao: 'WIII', name: 'Soekarno-Hatta International', city: 'Jakarta', country: 'Indonesia', latitude: -6.1256, longitude: 106.6558 },
  { iata: 'DEL', icao: 'VIDP', name: 'Indira Gandhi International', city: 'New Delhi', country: 'India', latitude: 28.5562, longitude: 77.1000 },
  { iata: 'BOM', icao: 'VABB', name: 'Chhatrapati Shivaji International', city: 'Mumbai', country: 'India', latitude: 19.0896, longitude: 72.8656 },
  { iata: 'SYD', icao: 'YSSY', name: 'Sydney Kingsford Smith', city: 'Sydney', country: 'Australia', latitude: -33.9399, longitude: 151.1753 },
  { iata: 'MEL', icao: 'YMML', name: 'Melbourne Airport', city: 'Melbourne', country: 'Australia', latitude: -37.6690, longitude: 144.8410 },
  { iata: 'BNE', icao: 'YBBN', name: 'Brisbane Airport', city: 'Brisbane', country: 'Australia', latitude: -27.3842, longitude: 153.1175 },
  { iata: 'AKL', icao: 'NZAA', name: 'Auckland Airport', city: 'Auckland', country: 'New Zealand', latitude: -37.0082, longitude: 174.7850 },
  { iata: 'GRU', icao: 'SBGR', name: 'São Paulo–Guarulhos', city: 'São Paulo', country: 'Brazil', latitude: -23.4356, longitude: -46.4731 },
  { iata: 'GIG', icao: 'SBGL', name: 'Rio de Janeiro–Galeão', city: 'Rio de Janeiro', country: 'Brazil', latitude: -22.8099, longitude: -43.2505 },
  { iata: 'EZE', icao: 'SAEZ', name: 'Ministro Pistarini', city: 'Buenos Aires', country: 'Argentina', latitude: -34.8222, longitude: -58.5358 },
  { iata: 'SCL', icao: 'SCEL', name: 'Arturo Merino Benítez', city: 'Santiago', country: 'Chile', latitude: -33.3930, longitude: -70.7858 },
  { iata: 'BOG', icao: 'SKBO', name: 'El Dorado International', city: 'Bogotá', country: 'Colombia', latitude: 4.7016, longitude: -74.1469 },
  { iata: 'JNB', icao: 'FAOR', name: 'O.R. Tambo International', city: 'Johannesburg', country: 'South Africa', latitude: -26.1367, longitude: 28.2411 },
  { iata: 'CPT', icao: 'FACT', name: 'Cape Town International', city: 'Cape Town', country: 'South Africa', latitude: -33.9715, longitude: 18.6021 },
  { iata: 'CAI', icao: 'HECA', name: 'Cairo International', city: 'Cairo', country: 'Egypt', latitude: 30.1219, longitude: 31.4056 },
  { iata: 'CMN', icao: 'GMMN', name: 'Mohammed V International', city: 'Casablanca', country: 'Morocco', latitude: 33.3675, longitude: -7.5898 },
  { iata: 'NBO', icao: 'HKJK', name: 'Jomo Kenyatta International', city: 'Nairobi', country: 'Kenya', latitude: -1.3192, longitude: 36.9278 },
]

export function searchAirports(query: string, limit = 10): Airport[] {
  if (!query || query.length < 2) return []
  
  const lowerQuery = query.toLowerCase()
  
  return airports
    .filter(airport => 
      airport.iata.toLowerCase().includes(lowerQuery) ||
      airport.name.toLowerCase().includes(lowerQuery) ||
      airport.city.toLowerCase().includes(lowerQuery) ||
      airport.country.toLowerCase().includes(lowerQuery)
    )
    .slice(0, limit)
}

export function getAirportByIata(iata: string): Airport | undefined {
  return airports.find(a => a.iata.toUpperCase() === iata.toUpperCase())
}


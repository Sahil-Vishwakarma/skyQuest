import { useEffect, useRef } from 'react'
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet'
import L from 'leaflet'
import type { Flight } from '../types'
import { airports, getAirportByIata } from '../utils/airports'

// Fix for default markers
import 'leaflet/dist/leaflet.css'

// Create custom aircraft icon
const createAircraftIcon = (rotation: number, isSelected: boolean) => {
  const color = isSelected ? '#fbbf24' : '#38bdf8'
  const size = isSelected ? 32 : 24
  
  return L.divIcon({
    html: `
      <svg width="${size}" height="${size}" viewBox="0 0 24 24" fill="${color}" style="transform: rotate(${rotation}deg); filter: drop-shadow(0 2px 4px rgba(0,0,0,0.5));">
        <path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/>
      </svg>
    `,
    className: 'aircraft-marker-icon',
    iconSize: [size, size],
    iconAnchor: [size / 2, size / 2],
  })
}

// Create airport icon
const createAirportIcon = (isSelected: boolean) => {
  const color = isSelected ? '#fbbf24' : '#94a3b8'
  const size = isSelected ? 16 : 10
  
  return L.divIcon({
    html: `
      <div style="
        width: ${size}px;
        height: ${size}px;
        background: ${color};
        border-radius: 50%;
        border: 2px solid rgba(255,255,255,0.5);
        box-shadow: 0 2px 4px rgba(0,0,0,0.3);
      "></div>
    `,
    className: 'airport-marker-icon',
    iconSize: [size, size],
    iconAnchor: [size / 2, size / 2],
  })
}

// Component to handle map updates
function MapUpdater({ flight }: { flight: Flight | null }) {
  const map = useMap()
  
  useEffect(() => {
    if (flight) {
      map.flyTo([flight.latitude, flight.longitude], 5, {
        duration: 1.5,
      })
    }
  }, [flight, map])
  
  return null
}

interface MapComponentProps {
  flight: Flight | null
  onAirportSelect: (iata: string | null) => void
  selectedAirport: string | null
}

export function MapComponent({ flight, onAirportSelect, selectedAirport }: MapComponentProps) {
  const mapRef = useRef<L.Map>(null)
  
  const selectedAirportData = selectedAirport ? getAirportByIata(selectedAirport) : null

  return (
    <div className="w-full h-full min-h-[400px] lg:min-h-0">
      <MapContainer
        center={flight ? [flight.latitude, flight.longitude] : [40, 0]}
        zoom={4}
        className="w-full h-full"
        ref={mapRef}
        zoomControl={true}
        attributionControl={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://carto.com/">CARTO</a>'
          url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
        />
        
        <MapUpdater flight={flight} />
        
        {/* Flight marker */}
        {flight && (
          <Marker
            position={[flight.latitude, flight.longitude]}
            icon={createAircraftIcon(flight.direction, true)}
          >
            <Popup>
              <div className="text-midnight-900">
                <p className="font-bold">{flight.flightNumber || flight.callsign || 'Unknown Flight'}</p>
                {flight.airline.name && <p className="text-sm">{flight.airline.name}</p>}
                <p className="text-sm text-midnight-600">
                  Alt: {Math.round(flight.altitude).toLocaleString()} ft
                </p>
                <p className="text-sm text-midnight-600">
                  Speed: {Math.round(flight.speed)} kts
                </p>
              </div>
            </Popup>
          </Marker>
        )}
        
        {/* Airport markers */}
        {airports.map((airport) => (
          <Marker
            key={airport.iata}
            position={[airport.latitude, airport.longitude]}
            icon={createAirportIcon(selectedAirport === airport.iata)}
            eventHandlers={{
              click: () => onAirportSelect(airport.iata),
            }}
          >
            <Popup>
              <div className="text-midnight-900">
                <p className="font-bold">{airport.iata}</p>
                <p className="text-sm">{airport.name}</p>
                <p className="text-xs text-midnight-600">{airport.city}, {airport.country}</p>
                <button
                  onClick={() => onAirportSelect(airport.iata)}
                  className="mt-2 px-3 py-1 bg-sky-500 text-white text-sm rounded-lg hover:bg-sky-600 transition-colors"
                >
                  Select as Guess
                </button>
              </div>
            </Popup>
          </Marker>
        ))}

        {/* Highlight selected airport */}
        {selectedAirportData && (
          <Marker
            position={[selectedAirportData.latitude, selectedAirportData.longitude]}
            icon={createAirportIcon(true)}
          />
        )}
      </MapContainer>
    </div>
  )
}


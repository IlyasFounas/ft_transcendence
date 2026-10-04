import 'leaflet/dist/leaflet.css';

import { useEffect } from 'react';
import L from 'leaflet';

type Props = {
    latitude: number,
    longitude: number,
    onGuess?: (lat: number, lng: number) => void,
    initialGuess?: { lat: number, lng: number } | null,
}

export default function Load2DMap({latitude, longitude, initialGuess, onGuess}: Props) {
  useEffect(() => {
    // Initializes the map and sets the view to the specified latitude and longitude with a zoom level
    const map = L.map('leaflet-map').setView([latitude, longitude], 2.5);

    // Adds the OpenStreetMap tile layer to the map with attribution and maximum zoom level
    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      minZoom: 2.5,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

        const marker = L.marker([latitude, longitude]);

        if (initialGuess) {
            marker.setLatLng([initialGuess.lat, initialGuess.lng]).addTo(map);
        }

        const onMapClick = (e: L.LeafletMouseEvent) => {
            marker
                .setLatLng(e.latlng)
                .addTo(map);
            onGuess?.(e.latlng.lat, e.latlng.lng);
        };

        map.on('click', onMapClick);

        return () => {
            map.remove();
        };
    }, [latitude, longitude]);
    return (
        <div id="leaflet-map" style={{ width: '100%', height: '100%' }}></div>
    );
}
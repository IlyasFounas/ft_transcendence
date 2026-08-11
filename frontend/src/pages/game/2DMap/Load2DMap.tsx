import 'leaflet/dist/leaflet.css';

import { useEffect } from 'react';
import L from 'leaflet';

type Props = {
  latitude: number,
  longitude: number,
}

export default function Load2DMap({latitude, longitude}: Props) {
  useEffect(() => {
    // Initializes the map and sets the view to the specified latitude and longitude with a zoom level
    const map = L.map('leaflet-map').setView([latitude, longitude], 3);

    // Adds the OpenStreetMap tile layer to the map with attribution and maximum zoom level
    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    // Creates a pin that can be moved around the map when the user clicks on a new location
    const marker = L.marker([latitude, longitude]);

    const onMapClick = (e: L.LeafletMouseEvent) => {
        marker
            .setLatLng(e.latlng)
            .addTo(map);
    };

    // Adds a click event listener
    map.on('click', onMapClick);
    onMapClick({ latlng: L.latLng(latitude, longitude) } as L.LeafletMouseEvent);

    return () => {
      map.remove();
    };
  }, [latitude, longitude]);
  return (
    <div id="leaflet-map" style={{ width: '100%', height: '100%' }}></div>
  );
}
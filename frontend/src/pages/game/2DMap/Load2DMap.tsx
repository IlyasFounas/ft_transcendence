import 'leaflet/dist/leaflet.css';

import { useEffect } from 'react';
import L from 'leaflet';

type Props = {
  latitude: number,
  longitude: number,
}

export default function Load2DMap({latitude, longitude}: Props) {
  useEffect(() => {
    const map = L.map('leaflet-map').setView([latitude, longitude], 13);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 3,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    return () => {
      map.remove();
    };
  }, [latitude, longitude]);
  return (
    <div id="leaflet-map" style={{ width: '100%', height: '100%' }}></div>
  );
}
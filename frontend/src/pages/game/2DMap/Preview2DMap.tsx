import 'leaflet/dist/leaflet.css';

import { useEffect } from 'react';
import L from 'leaflet';


const WORLD_CENTER: L.LatLngTuple = [20, 0];

export default function Preview2DMap() {
  useEffect(() => {
    const map = L.map('leaflet-map', {
      zoomControl: false,
      dragging: false,
      scrollWheelZoom: false,
      doubleClickZoom: false,
      touchZoom: false,
      boxZoom: false,
      keyboard: false,
    }).setView(WORLD_CENTER, 3);

    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    return () => {
      map.remove();
    };
  }, []);

  return (
    <div
      id="leaflet-map"
      style={{ width: '100%', height: '100%' }}
    />
  );
}
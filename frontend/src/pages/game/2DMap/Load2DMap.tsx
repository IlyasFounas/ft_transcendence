import 'leaflet/dist/leaflet.css';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

import { useEffect } from 'react';
import L from 'leaflet';

type Props = {
  onGuess?: (lat: number, lng: number) => void,
  initialGuess?: { lat: number, lng: number } | null,
  initialView?: { center: L.LatLngTuple, zoom: number } | null,
  onViewChange?: (center: L.LatLngTuple, zoom: number) => void,
}

const WORLD_CENTER: L.LatLngTuple = [20, 0];

export default function Load2DMap({initialGuess, onGuess, initialView, onViewChange}: Props) {
  useEffect(() => {

	const map = L.map('leaflet-map').setView(
		initialView?.center ?? WORLD_CENTER,
		initialView?.zoom ?? 2.5
		);
	// Adds the OpenStreetMap tile layer to the map with attribution and maximum zoom level
	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19,
		minZoom: 2.5,
		attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);
	
	   const icon = L.icon({
			iconUrl: markerIcon,
			iconRetinaUrl: markerIcon2x,
			shadowUrl: markerShadow,
			iconSize: [25, 41],
			iconAnchor: [12, 41],
			shadowSize: [41, 41],
		});
		const marker = L.marker(WORLD_CENTER, {icon: icon});

		if (initialGuess) {
			marker.setLatLng([initialGuess.lat, initialGuess.lng]).addTo(map);
		}

		const onMapClick = (e: L.LeafletMouseEvent) => {
			marker
				.setLatLng(e.latlng)
				.addTo(map);
			onGuess?.(e.latlng.lat, e.latlng.lng);
		};

		map.on('moveend', () => {
			const c = map.getCenter();
			onViewChange?.([c.lat, c.lng], map.getZoom());
			});

		return () => {
			map.remove();
		};
	}, []);
	return (
		<div id="leaflet-map" style={{ width: '100%', height: '100%' }}></div>
	);
}
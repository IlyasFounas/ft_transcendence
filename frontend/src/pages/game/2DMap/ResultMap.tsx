import 'leaflet/dist/leaflet.css';
import { renderToStaticMarkup } from "react-dom/server";

import { useEffect, useRef } from 'react';
import L from 'leaflet';
import Avatar from '../../../components/primitives/Avatar';

type Location = {
  latitude: number;
  longitude: number;
}

type PlayerGuesses = {
	location: Location;
	username: string;
	avatar: string;
}

type Props = {
	actualLocation: Location;
	guesses : PlayerGuesses[];
	currentUsername: string;
	hoveredPlayerId?: string | null;
}

function createDotIcon() {
  return L.divIcon({
	className: "",
	html: `
	  <div class="relative w-6 h-6 flex items-center justify-center">
		<div class="w-1.5 h-1.5 rounded-full bg-slate-700"></div>
	  </div>
	`,
	iconSize: [24, 24],
	iconAnchor: [12, 12],
  });
}

function createTargetIcon(color: string) {
  return L.divIcon({
	className: "",
	html: `
	  <div class="relative w-8 h-8 flex items-center justify-center">
		<div class="absolute w-4 h-4 rounded-full opacity-60 animate-ping" style="background:${color}"></div>
		<svg viewBox="0 0 32 32" class="relative w-8 h-8">
		  <circle cx="16" cy="16" r="10" fill="none" stroke="${color}" strokeWidth="2"/>
		  <circle cx="16" cy="16" r="3" fill="${color}"/>
		  <line x1="16" y1="0" x2="16" y2="6" stroke="${color}" strokeWidth="2"/>
		  <line x1="16" y1="26" x2="16" y2="32" stroke="${color}" strokeWidth="2"/>
		  <line x1="0" y1="16" x2="6" y2="16" stroke="${color}" strokeWidth="2"/>
		  <line x1="26" y1="16" x2="32" y2="16" stroke="${color}" strokeWidth="2"/>
		</svg>
	  </div>
	`,
	iconSize: [32, 32],
	iconAnchor: [16, 16],
  });
}

function PlayerTooltip({ avatar, username }: { avatar: string, username: string})
{
	return (
		<div className="flex flex-row  justify-center items-center gap-2">
			<Avatar avatar={avatar} size={35}></Avatar>
			<span>{username}</span>
		</div>
	)
}

export default function ResultMap({actualLocation, guesses, currentUsername, hoveredPlayerId}: Props) {
  
	const guessMarkersRef = useRef<Map<string, L.Marker>>(new Map());

	useEffect(() => {
		const map = L.map('leaflet-result-map', {minZoom: 2});

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
	  maxZoom: 19,
	  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	L.marker([actualLocation.latitude, actualLocation.longitude], {icon: createTargetIcon("#22c55e")}).addTo(map);

	guessMarkersRef.current.clear();

	guesses.forEach((guess) => {
  		const isCurrentPlayer = currentUsername === guess.username;

		const marker = L.marker(
			[guess.location.latitude, guess.location.longitude],
			{ icon: isCurrentPlayer ? createTargetIcon("#ef4444") : createDotIcon() }
		)
		.addTo(map)
		.bindTooltip(
			renderToStaticMarkup(<PlayerTooltip avatar={guess.avatar} username={guess.username} />),
			{ direction: "top", offset: [0, -12], className: "player-tooltip" }
		);

		guessMarkersRef.current.set(guess.username, marker);

		L.polyline(
			[
			[actualLocation.latitude, actualLocation.longitude],
			[guess.location.latitude, guess.location.longitude],
			],
			{
			color:"#000000", 
			weight: 1,
			dashArray: "6, 3"
			}
		).addTo(map);
		});

	const currentPlayerGuess = guesses.find((g) => g.username === currentUsername);

	const bounds = L.latLngBounds([
	[actualLocation.latitude, actualLocation.longitude],
	...(currentPlayerGuess
		? [[currentPlayerGuess.location.latitude, currentPlayerGuess.location.longitude] as [number, number]]
		: []),
	]);

	map.fitBounds(bounds, { padding: [40, 40] });

	return () => {
	 	 map.remove();
		};
  	}, [guesses, actualLocation.latitude, actualLocation.longitude]);

	useEffect(() => {
		guessMarkersRef.current.forEach((marker, username) => {
			const el = marker.getElement();
			if (!el) return;
			const isHovered= username === hoveredPlayerId;

			if (isHovered) {
				marker.openTooltip();
			} else {
				marker.closeTooltip();
			}
		});
		}, [hoveredPlayerId]);

  return (
	<div id="leaflet-result-map" style={{ width: '100%', height: '100%' }}></div>
  );
}
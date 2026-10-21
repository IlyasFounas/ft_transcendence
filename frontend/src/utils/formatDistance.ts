export function formatDistance(meters: number): { value: number; unit: string } {
	if (meters < 1000) {
		return { value: Math.round(meters), unit: "m" };
	}
	return { value: Math.round(meters / 1000), unit: "km" };
}
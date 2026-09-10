import type { Game } from "../types/game";

export const mockGame: Game = {
    id: 42,
    roomcode: "214748",
    currentRound: 8,
    totalRounds: 5,

    totalTime: 12,

    location: {
        latitude: 45.782429,
        longitude: 4.745860,
    },

    streetView: {
        id: "271688394642360",
        moving: true,
        zooming: true,
        panning: true,
    },
};

// Lieux réels pour chaque round (à indexer par numéro de round)
export const mockRoundLocations: Record<number, { latitude: number; longitude: number }> = {
    1: { latitude: 48.8566, longitude: 2.3522 },   // Paris
    2: { latitude: 41.9028, longitude: 12.4964 },  // Rome
    3: { latitude: 52.5200, longitude: 13.4050 },  // Berlin
    4: { latitude: 40.4168, longitude: -3.7038 },  // Madrid
    5: { latitude: 59.3293, longitude: 18.0686 },  // Stockholm
};
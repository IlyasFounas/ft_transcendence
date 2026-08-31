import type { Game } from "../types/game";

export const mockGame: Game = {
    id: 42,

    currentRound: 8,
    totalRounds: 10,

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
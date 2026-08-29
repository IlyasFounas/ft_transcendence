export interface Player {
    username: string;
    role: string;
    avatar: string;
}

export interface PlayerResult {
    player: Player;
    rounds: RoundResult[];
}

export interface GameResult {
    currentRound: number;
    totalRounds: number;
    players: PlayerResult[];
    currentUsername: string;
    distance: number;
}

export type RoundResult = {
  round: number;
  points: number;
  distance: number;
  guessLocation: {
    latitude: number;
    longitude: number;
  };
};

export interface Game {
    id: number;
    currentRound: number;
    totalRounds: number;
    totalTime: number;

    location: {
        latitude: number;
        longitude: number;
    };

    streetView: {
        id: string;
        moving: boolean;
        zooming: boolean;
        panning: boolean;
    };
}
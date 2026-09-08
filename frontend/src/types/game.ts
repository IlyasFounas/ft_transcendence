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
	roomcode: string;
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

export interface GameRoom {
    id: number;
    title: string;
    mode: string;
    difficulty: string;
    timer: number;
    roundNb: number;
    finished: boolean;
    playerNb: number; // nombre max de joueurs
    createdAt: string;
}

export interface Lobby {
    id: number;
    userId: number;
    gameId: number;
    joinCode: string | null;
    isCreator: boolean;
    leftAt: string | null;
}
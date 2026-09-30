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
    playerNb: number;
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

export interface LocationInfo {
	lat: number;
	long: number;
	country_code: string;
}

export interface RoundInfos {
	distance: number;
	score: number;
	base_score: number;
	country_found: boolean;
	base: LocationInfo;
	guess: LocationInfo;
	streak: number;
	time_left_seconds: number;
}

export interface PlayerRoundResult {
	user_id: number;
	username: string;
	round_nb: number;
	has_guessed: boolean;
	round_infos: RoundInfos | null;
	total_score: number;
}
export interface Player {
  username: string;
  role: string;
  avatar: string;
}

export interface PlayerResult {
  rank: number;
  player: Player;
  rounds: RoundResult[];
}

export interface GameResult {
  currentRound: number;
  totalRounds: number;
  players: PlayerResult[];
  currentPlayerUsername: string; 
  distance: number;
}

export interface RoundResult {
  round: number;
  points: number;
  distance: number;
}
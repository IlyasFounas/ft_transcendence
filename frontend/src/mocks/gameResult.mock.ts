import type { GameResult, RoundResult } from "../types/game";

function generateRounds(count: number): RoundResult[] {
  return Array.from({ length: count }, (_, i) => ({
    round: i + 1,
    points: 800 + Math.floor(Math.random() * 700),
    distance: Math.floor(Math.random() * 500),
  }));
}

export const mockGameResult: GameResult = {
  currentRound: 100,
  totalRounds: 100,
  currentPlayerUsername: "Lucie",
  distance: 4225,
  players: Array.from({ length: 10 }, (_, i) => ({
    rank: i + 1,
    player: {
      username: i === 0 ? "Lucie" : `Player${i + 1}`,
      role: "Player",
      avatar: "/avatars/avatar1.png",
    },
    rounds: generateRounds(100),
  })),
};


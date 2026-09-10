import type { GameResult, RoundResult } from "../types/game";

function generateRounds(count: number): RoundResult[] {
  return Array.from({ length: count }, (_, i) => ({
    round: i + 1,
    points: 800 + Math.floor(Math.random() * 700),
    distance: Math.floor(Math.random() * 500),
    guessLocation: {
      // coordonnées aléatoires quelque part en Europe, juste pour tester visuellement
      latitude: 35 + Math.random() * 25,   // ~35 à 60
      longitude: -10 + Math.random() * 40, // ~-10 à 30
    },
  }));
}

export const mockGameResult: GameResult = {
  currentRound: 5,
  totalRounds: 5,
  currentUsername: "Lucie",
  distance: 4225,

  players: Array.from({ length: 7 }, (_, i) => ({
    player: {
      username: i === 0 ? "Lucie" : `Player${i + 1}`,
      role: "Player",
      avatar: "/avatars/avatar1.png",
    },

    rounds: generateRounds(5),
  })),
};
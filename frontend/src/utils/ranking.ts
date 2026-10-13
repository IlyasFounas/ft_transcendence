import type { PlayerResult } from "../types/game";

export function getSortedPlayersWithRank(players: PlayerResult[]) {
  return [...players]
	.sort((a, b) => {
	  const totalA = a.rounds.reduce((sum, r) => sum + r.points, 0);
	  const totalB = b.rounds.reduce((sum, r) => sum + r.points, 0);
	  return totalB - totalA;
	})
	.map((result, index) => ({ ...result, rank: index + 1 }));
}
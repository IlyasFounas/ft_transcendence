import type { RoundResult } from "../types/game";

export function getTotalPoints(rounds: RoundResult[]): number {
	return rounds.reduce((total, r) => total + r.points, 0);
}
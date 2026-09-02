import type { GameResult } from "../types/game";
import { useMemo } from "react";

export function useRoundGuesses(gameResult: GameResult) {
	const guesses = useMemo(
		() =>
		  gameResult.players
			.map((p) => {
			  const round = p.rounds.find((r) => r.round === gameResult.currentRound);
			  if (!round) return null;
			  return {
				username: p.player.username,
				location: round.guessLocation,
				avatar: p.player.avatar,
			  };
			})
			.filter((g): g is NonNullable<typeof g> => g !== null),
		[gameResult.players, gameResult.currentRound]
	  );

	return (guesses);
}
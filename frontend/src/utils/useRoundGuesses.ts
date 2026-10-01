import type { GameResult } from "../types/game";
import { useMemo } from "react";

export function useRoundGuesses(gameResult: GameResult, round?: number) {
	const targetRound = round ?? gameResult.currentRound;
	
	const guesses = useMemo(
		() =>
			gameResult.players
				.map((p) => {
					const roundData = p.rounds.find((r) => r.round === targetRound);
					if (!roundData || (roundData.guessLocation.latitude === 0 && roundData.guessLocation.longitude === 0)) {
                        return null;
                    }
					return {
						username: p.player.username,
						location: roundData.guessLocation,
						avatar: p.player.avatar,
					};
				})
				.filter((g): g is NonNullable<typeof g> => g !== null),
		[gameResult.players, targetRound]
	);

	return guesses;
}
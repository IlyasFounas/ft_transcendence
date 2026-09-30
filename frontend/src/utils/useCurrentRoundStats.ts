import { useMemo } from "react";
import type { GameResult } from "../types/game";
import { useCountUp } from "./useCountUp";
import { formatDistance } from "./formatDistance";

export function useCurrentRoundStats(gameResult: GameResult) {
	const currentPlayerResult = useMemo(
		() => gameResult.players.find((p) => p.player.username === gameResult.currentUsername),
		[gameResult.players, gameResult.currentUsername]
	);

	const currentRoundResult = useMemo(() => {
		return currentPlayerResult?.rounds.find((r) => r.round === gameResult.currentRound);
	}, [currentPlayerResult, gameResult.currentRound]);

	const currentRoundScore = currentRoundResult?.points ?? 0;
	const totalScore = currentPlayerResult?.rounds.reduce((sum, r) => sum + r.points, 0) ?? 0;

	const { displayValue: animatedScore, justFinished: scoreJustFinished } = useCountUp(currentRoundScore);

	const currentRoundDistance = currentRoundResult?.distance ?? 0;
	const { value: distanceValue, unit: distanceUnit } = formatDistance(currentRoundDistance);
	const { displayValue: animatedDistance, justFinished: distanceJustFinished } = useCountUp(distanceValue);

	const totalDistance = currentPlayerResult?.rounds.reduce((sum, r) => sum + r.distance, 0) ?? 0;

	return {
		totalScore,
		currentPlayerResult,
		currentRoundDistance,
		currentRoundResult,
		animatedScore,
		scoreJustFinished,
		animatedDistance,
		distanceUnit,
		distanceJustFinished,
		totalDistance
	};
}
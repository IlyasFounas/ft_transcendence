import { useMemo } from "react";

import NextRoundTimer from "../../components/results/NextRoundTimer";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import RankResultsList from "../../components/results/RankResultsList";
import ResultMap from "./2DMap/ResultMap";
import { mockGameResult } from "../../mocks/gameResult.mock";
import { mockGame } from "../../mocks/game.mock";
import { useCountUp } from "../../utils/useCountUp";

export default function MultiplayerResultsPage() {
  const gameResult = mockGameResult;

  const currentRoundResult = useMemo(() => {
	const currentPlayerResult = gameResult.players.find(
	  (p) => p.player.username === gameResult.currentUsername
	);
	return currentPlayerResult?.rounds.find(
	  (r) => r.round === gameResult.currentRound
	);
  }, [gameResult.players, gameResult.currentUsername, gameResult.currentRound]);

  const currentRoundScore = currentRoundResult?.points ?? 0;

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
		.filter(
		  (g): g is { username: string; location: typeof g.location; avatar: string } =>
			g !== null
		),
	[gameResult.players, gameResult.currentRound]
  );

  const { displayValue: animatedScore, justFinished: scoreJustFinished } =
	useCountUp(currentRoundScore);
  const { displayValue: animatedDistance, justFinished: distanceJustFinished } =
	useCountUp(gameResult.distance);

  return (
	<div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px] gap-4 p-4">
	  <div className="order-1 lg:col-start-1 lg:row-start-1">
		<ResultInfoRow
		  title="Mission Debrief"
		  infos={`Round ${gameResult.currentRound} of ${gameResult.totalRounds}`}
		  variant="primary"
		/>
	  </div>

	  <div
		className="order-2 lg:col-start-1 lg:row-start-2 w-full h-60 lg:h-[650px]
		  outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden"
	  >
		<ResultMap
		  actualLocation={mockGame.location}
		  guesses={guesses}
		  currentUsername={gameResult.currentUsername}
		/>
	  </div>

	  <div className="flex flex-row order-3 gap-2 lg:col-start-1 lg:row-start-3">
		<ResultInfoRow
		  title="Score"
		  infos={String(animatedScore)}
		  variant="secondary"
		  suffix="pts"
		  animate={scoreJustFinished}
		/>
		<ResultInfoRow
		  title="Distance"
		  infos={String(animatedDistance)}
		  variant="secondary"
		  suffix="km"
		  animate={distanceJustFinished}
		/>
	  </div>

	  <div className="order-4 lg:col-start-2 lg:row-span-3 flex flex-col items-center gap-4">
		<ResultInfoRow title="Round Ranking" variant="primary" />
		<RankResultsList
		  players={gameResult.players}
		  currentRound={gameResult.currentRound}
		  currentUsername={gameResult.currentUsername}
		/>
		<NextRoundTimer />
	  </div>
	</div>
  );
}
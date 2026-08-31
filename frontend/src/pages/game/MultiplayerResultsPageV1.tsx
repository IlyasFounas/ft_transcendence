import GameResultList from "../../components/results/GameResultList";
import NextRoundTimer from "../../components/results/NextRoundTimer";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import RankResultsList from "../../components/results/RankResultsList";
import TotalScoreCard from "../../components/results/TotalScoreCard";
import ResultMap from "./2DMap/ResultMap";
import { mockGameResult } from "../../mocks/gameResult.mock";
import { mockGame } from "../../mocks/game.mock";

export default function MultiplayerResultsPage() {
  const gameResult = mockGameResult;

  const currentPlayerResult = gameResult.players.find(
	(p) => p.player.username === gameResult.currentUsername
  );
  const totalScore =
	currentPlayerResult?.rounds.reduce((sum, r) => sum + r.points, 0) ?? 0;

  const guesses = gameResult.players
	.map((p) => {
	  const round = p.rounds.find((r) => r.round === gameResult.currentRound);
	  if (!round) return null;
	  return {
		username: p.player.username,
		location: round.guessLocation,
		avatar: p.player.avatar,
	  };
	})
	.filter((g): g is NonNullable<typeof g> => g !== null)

  return (
	<div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px] gap-4 p-4">
	  <div className="order-1 lg:col-start-1 lg:row-start-1">
		<ResultInfoRow
		  title="Mission Debrief"
		  infos={`Round ${gameResult.currentRound} of ${gameResult.totalRounds}`}
		  variant="primary"
		/>
	  </div>

	  <div className="order-2 lg:col-start-1 lg:row-start-2 w-full h-60 lg:h-[380px] outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden">
		<ResultMap
		  actualLocation={mockGame.location}
		  guesses={guesses}
		  currentUsername={gameResult.currentUsername}
		/>
	  </div>

	  <div className="flex flex-row order-3 gap-2 lg:hidden">
		<ResultInfoRow
		  title="Total Score"
		  infos={String(totalScore)}
		  variant="secondary"
		  suffix="pts"
		/>
		<ResultInfoRow
		  title="Distance"
		  infos={String(gameResult.distance)}
		  variant="secondary"
		  suffix="km"
		/>
	  </div>

	  <div className="hidden lg:flex lg:col-start-2 lg:row-start-1">
		<ResultInfoRow title="Your results" variant="primary" />
	  </div>

	  <div className="order-4">
		<RankResultsList
		  players={gameResult.players}
		  currentRound={gameResult.currentRound}
		  currentUsername={gameResult.currentUsername}
		/>
	  </div>

	  <div className="order-5 lg:col-start-2 lg:row-start-2 lg:row-span-3 flex flex-col items-center gap-4">
		<div className="hidden lg:flex w-full">
		  <GameResultList rounds={currentPlayerResult?.rounds ?? []} currentRound={3}/>
		</div>
		<TotalScoreCard points={totalScore} distance={gameResult.distance} />
		<NextRoundTimer />
	  </div>
	</div>
  );
}
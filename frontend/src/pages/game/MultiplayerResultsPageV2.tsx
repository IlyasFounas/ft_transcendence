import { useState } from "react"

import NextRoundTimer from "../../components/results/NextRoundTimer";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import RankResultsList from "../../components/results/RankResultsList";
import ResultMap from "./2DMap/ResultMap";
import { mockGameResult } from "../../mocks/gameResult.mock";
import { mockGame } from "../../mocks/game.mock";
import { useCurrentRoundStats } from "../../utils/useCurrentRoundStats";
import { useRoundGuesses } from "../../utils/useRoundGuesses";

export default function MultiplayerResultsPage() {

  	const gameResult = mockGameResult;
  	const {
		animatedScore,
	 	scoreJustFinished,
	 	animatedDistance,
	 	distanceJustFinished,
	} = useCurrentRoundStats(gameResult);

	const guesses = useRoundGuesses(gameResult);
	const [hoveredPlayerId, setHoveredPlayerId] = useState<string | null>(null);

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
			hoveredPlayerId={hoveredPlayerId}
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
				onHoverPlayer={setHoveredPlayerId}
			/>
			<NextRoundTimer />
		</div>
		</div>
  );
}
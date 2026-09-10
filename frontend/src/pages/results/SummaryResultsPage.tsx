import RankResultsList from "../../components/results/RankResultsList"
import { mockGameResult } from "../../mocks/gameResult.mock";
import ResultMap from "../game/2DMap/ResultMap";
import { useRoundGuesses } from "../../utils/useRoundGuesses";
import { mockGame } from "../../mocks/game.mock";
import { mockRoundLocations } from "../../mocks/game.mock";
import { useState, useMemo } from "react"
import PrimaryButton from "../../components/buttons/PrimaryButton";
import {ChevronLeft, ChevronRight} from "lucide-react"
import ResultInfoRow from "../../components/results/ResultInfoRow";

export default function SummaryResultsPage()
{
  	const gameResult = mockGameResult;
	const playerCount = mockGameResult.players.length;
	const columns = playerCount > 5 ? 2 : 1;
	const guesses = useRoundGuesses(gameResult);
	const [hoveredPlayerId, setHoveredPlayerId] = useState<string | null>(null);
	const [roundSelected, setRoundSelected] = useState(1);

	const currentPlayerResult = gameResult.players.find(
	(p) => p.player.username === gameResult.currentUsername
	);
	const selectedRoundData = currentPlayerResult?.rounds.find(
		(r) => r.round === roundSelected
	);

	const roundRanking = useMemo(() => {
    const pointsForRound = gameResult.players
        .map((p) => ({
            username: p.player.username,
            points: p.rounds.find((r) => r.round === roundSelected)?.points ?? 0,
        }))
        .sort((a, b) => b.points - a.points);

    const position = pointsForRound.findIndex(
        (p) => p.username === gameResult.currentUsername
    );

    return position === -1 ? null : position + 1;
	}, [gameResult.players, gameResult.currentUsername, roundSelected]);


	return (
		<div className="w-full h-[calc(100vh-5rem)] p-5 flex flex-col gap-5
		bg-[url(/backgrounds/earth_background.png)]">
			<ResultInfoRow title="Global ranking" variant="primary"></ResultInfoRow>
			<div className="shrink-0">
				<RankResultsList
					players={mockGameResult.players}
					currentUsername={mockGameResult.currentUsername}
					columns={columns}
					onHoverPlayer={setHoveredPlayerId}
				/>
			</div>
			<div className="flex flex-col lg:flex-row gap-4 flex-1 lg:min-h-0">
    			<div className="w-full bg-surface-base outline outline-1 outline-sky-950 rounded-[10px] text-center overflow-hidden flex flex-col p-4 gap-4 border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] min-h-0">
					<div className="flex flex-row justify-between text-heading-medium uppercase">
						<button
						className={`
						transition-transform
						duration-200
						ease-out
						${roundSelected > 1 ? "text-white hover:scale-[1.2] active:scale-[0.98] cursor-pointer" : "text-gray-600"}
						`}
						onClick={() => {
							if (roundSelected > 1) {
								setRoundSelected(roundSelected - 1);
							}
						}}
						>
						<ChevronLeft />
					</button>
						<h1>Round {roundSelected}</h1>
						<button
						className={`
						transition-transform
						duration-200
						ease-out
						${roundSelected < mockGame.totalRounds ? "text-white hover:scale-[1.2] active:scale-[0.98] cursor-pointer" : "text-gray-600"}
						`}
						onClick={() => {
							if (roundSelected < mockGame.totalRounds) {
								setRoundSelected(roundSelected + 1);
							}
						}}
						>
						<ChevronRight />
					</button>
					</div>
					<div className="w-full h-[300px] lg:h-auto lg:flex-1 lg:min-h-0 outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden">
						<ResultMap
							actualLocation={mockRoundLocations[roundSelected]}
							guesses={guesses}
							currentUsername={gameResult.currentUsername}
							hoveredPlayerId={hoveredPlayerId}
						/>
					</div>
				</div>
				<div className="w-full lg:w-[400px] flex flex-col gap-4">
					<div className="w-full text-heading-small bg-surface-base uppercase outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden flex flex-col border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)]">
						<h1 className="uppercase flex p-4 text-text-primary/70 bg-linear-to-t from-black to-surface-base">
							Your results
						</h1>

						<div className="w-full h-[80px] outline outline-1 outline-sky-950 flex flex-row items-center justify-between px-4">
							<h1>Score</h1>
							<h1>{selectedRoundData?.points ?? "-"} pts</h1>
						</div>

						<div className="w-full h-[80px] outline outline-1 outline-sky-950 flex flex-row items-center justify-between px-4">
							<h1>Distance</h1>
							<h1>{selectedRoundData?.distance ?? "-"} km</h1>
						</div>

						<div className="w-full h-[80px] outline outline-1 outline-sky-950 flex flex-row items-center justify-between px-4">
							<h1>Round Ranking</h1>
							<h1>{roundRanking ?? "-"}/{gameResult.players.length}</h1>
						</div>
					</div>
					<div className="flex justify-center">
						<PrimaryButton size="medium">Return To Home</PrimaryButton>
					</div>
				</div>
			</div>
		</div>
	)
}
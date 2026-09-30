import PrimaryButton from "../../components/buttons/PrimaryButton";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import TotalScoreCard from "../../components/results/TotalScoreCard";
import GameResultList from "../../components/results/GameResultList";
import ResultMap from "../game/2DMap/ResultMap";
import { useCurrentRoundStats } from "../../utils/useCurrentRoundStats";
import { useRoundGuesses } from "../../utils/useRoundGuesses";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";
import { useLobbySocket } from "../../context/GameSocketContext";
import type { GameResult, RoundResult } from "../../types/game";

export default function SoloResultPage() {
	const location = useLocation();
	const navigate = useNavigate();
	const { user } = useAuth();
	const { startNextRound, leaveLobby, resetGame } = useLobbySocket();

	const { currentRound, totalRounds, rounds, actualLocation } = location.state as {
		currentRound: number;
		totalRounds: number;
		rounds: RoundResult[];
		actualLocation: { latitude: number; longitude: number };
	};

	const lastRound = rounds[rounds.length - 1];

	const gameResult: GameResult = {
		currentRound,
		totalRounds,
		currentUsername: user?.Username ?? "",
		distance: lastRound?.distance ?? 0,
		players: [
			{
				player: {
					username: user?.Username ?? "",
					role: "Player",
					avatar: user?.Avatar ?? "/avatars/avatar1.png",
				},
				rounds,
			},
		],
	};

	const {
		totalScore,
		currentPlayerResult,
		animatedScore,
		scoreJustFinished,
		animatedDistance,
		distanceJustFinished,
		distanceUnit,
		totalDistance,
	} = useCurrentRoundStats(gameResult);
	
	const guesses = useRoundGuesses(gameResult);

	const isLastRound = currentRound >= totalRounds;

	const handleNext = () => {
		if (isLastRound) {
			resetGame();
			leaveLobby();
			navigate("/home"); 
		} else {
			startNextRound();
			navigate("/game/solo"); 
		}
	};

	return (
		<div className="
			grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px]
			lg:grid-rows-[auto_1fr_auto]
			lg:h-[calc(100vh-2rem)]
			gap-4
			p-4">

			<div className="order-1 lg:col-start-1 lg:row-start-1">
				<ResultInfoRow
					title="Mission Debrief"
					infos={`Round ${gameResult.currentRound} of ${totalRounds}`}
					variant="primary"
				/>
			</div>

			<div className="
				order-2
				lg:col-start-1 lg:row-start-2
				w-full h-60 lg:h-full
				outline outline-1 outline-sky-950
				rounded-[10px]
				overflow-hidden
			">
				<ResultMap
					actualLocation={actualLocation}
					guesses={guesses}
					currentUsername={gameResult.currentUsername}
				/>
			</div>

			<div className="
				flex flex-row order-3 gap-2
				lg:col-start-1 lg:row-start-3
			">
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
					suffix={distanceUnit}
					animate={distanceJustFinished}
				/>
			</div>
			<div className="hidden lg:flex lg:col-start-2 lg:row-start-1">
				<ResultInfoRow title="Your results" variant="primary"/>
			</div>
			<div className="
				order-4
				lg:col-start-2 lg:row-start-2 lg:row-span-3
				flex flex-col
				items-center
				gap-4
			">
				<GameResultList currentRound={gameResult.currentRound} rounds={currentPlayerResult?.rounds ?? []} totalRounds={totalRounds} />
				<TotalScoreCard points={totalScore} distance={totalDistance} />
				
				<PrimaryButton size="medium" onClick={handleNext}>
					{isLastRound ? "Finish Mission" : "Next Round"}
				</PrimaryButton>
			</div>

		</div>
	)
}
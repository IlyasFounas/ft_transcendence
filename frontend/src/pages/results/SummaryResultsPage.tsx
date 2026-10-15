import { useLocation, useNavigate } from "react-router-dom";
import RankResultsList from "../../components/results/RankResultsList"
import ResultMap from "../game/2DMap/ResultMap";
import { useRoundGuesses } from "../../utils/useRoundGuesses";
import { useState, useMemo } from "react"
import PrimaryButton from "../../components/buttons/PrimaryButton";
import {ChevronLeft, ChevronRight} from "lucide-react"
import ResultInfoRow from "../../components/results/ResultInfoRow";
import { useAuth } from "../../context/AuthContext";
import { useLobbySocket } from "../../context/GameSocketContext";
import { formatDistance } from "../../utils/formatDistance";

export default function SummaryResultsPage()
{
	const location = useLocation();
	const navigate = useNavigate();
	const { user } = useAuth();
	const { leaveLobby, resetGame } = useLobbySocket();
	
	const { fullHistory, players, totalRounds } = location.state || {};

	const [hoveredPlayerId, setHoveredPlayerId] = useState<string | null>(null);
	const [roundSelected, setRoundSelected] = useState(1);

	const gameResult = useMemo(() => {
		if (!fullHistory || fullHistory.length === 0) return null;

		const safePlayers = (players && players.length > 0) 
			? players 
			: fullHistory[0].results.map((r: any) => ({
				id: r.user_id,
				username: r.username,
				avatar: "/avatars/avatar1.png",
				isHost: false
			}));

		const formattedPlayers = safePlayers.map((p: any) => {
			const playerRounds = fullHistory.map((round: any) => {
				const result = round.results.find((r: any) => r.user_id === p.id);
				return {
					round: round.round_nb,
					points: result?.round_infos?.score ?? 0,
					distance: result?.round_infos?.distance ?? 0,
					guessLocation: result?.has_guessed && result.round_infos ? {
						latitude: result.round_infos.guess.lat,
						longitude: result.round_infos.guess.long
					} : { latitude: 0, longitude: 0 }
				};
			});

			return {
				player: {
					username: p.username,
					role: p.isHost ? "Host" : "Player",
					avatar: p.avatar ?? "/avatars/avatar1.png",
					id: p.id.toString()
				},
				rounds: playerRounds
			};
		});

		const currentUsername = user?.Username ?? "";
		const currentPlayer = formattedPlayers.find((p: any) => p.player.username === currentUsername);
		const totalDistance = currentPlayer?.rounds.reduce((sum: number, r: any) => sum + r.distance, 0) ?? 0;

		return {
			currentRound: totalRounds,
			totalRounds: totalRounds,
			currentUsername: currentUsername,
			distance: totalDistance,
			players: formattedPlayers,
		};
	}, [fullHistory, players, totalRounds, user]);

	const roundLocations = useMemo(() => {
		if (!fullHistory) return {};
		const locations = {};
		fullHistory.forEach(round => {
			locations[round.round_nb] = {
				latitude: round.base.lat,
				longitude: round.base.long
			};
		});
		return locations;
	}, [fullHistory]);

	const roundRanking = useMemo(() => {
		if (!gameResult) return null;

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
	}, [gameResult, roundSelected]);

	const guesses = useRoundGuesses(gameResult, roundSelected);

	const handleReturnHome = () => {
		leaveLobby();
		resetGame();
		navigate("/");
	};

   if (!gameResult) {
		return (
			<div className="flex h-screen items-center justify-center text-white bg-slate-900">
				<p className="text-xl animate-pulse">Loading mission debrief...</p>
			</div>
		);
	}

	const playerCount = gameResult.players.length;
	const columns = playerCount > 5 ? 2 : 1;

	const currentPlayerResult = gameResult.players.find(
		(p) => p.player.username === gameResult.currentUsername
	);
	const selectedRoundData = currentPlayerResult?.rounds.find(
		(r) => r.round === roundSelected
	);

	const formattedDistance = selectedRoundData
	? formatDistance(selectedRoundData.distance)
	: null;

	return (
		<div className="w-full h-screen p-5 flex flex-col gap-5
		bg-[url(/backgrounds/earth_background.png)] bg-cover bg-no-repeat bg-center">
			<ResultInfoRow title="Global ranking" variant="primary"></ResultInfoRow>
			<div className="shrink-0">
				<RankResultsList
					players={gameResult.players}
					currentUsername={gameResult.currentUsername}
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
						${roundSelected < gameResult.totalRounds ? "text-white hover:scale-[1.2] active:scale-[0.98] cursor-pointer" : "text-gray-600"}
						`}
						onClick={() => {
							if (roundSelected < gameResult.totalRounds) {
								setRoundSelected(roundSelected + 1);
							}
						}}
						>
						<ChevronRight />
					</button>
					</div>
					<div className="w-full h-[300px] lg:h-auto lg:flex-1 lg:min-h-0 outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden">
						<ResultMap
							actualLocation={roundLocations[roundSelected]}
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
							<h1>
								{formattedDistance
									? `${formattedDistance.value} ${formattedDistance.unit}`
									: "-"}
							</h1>
						</div>

						<div className="w-full h-[80px] outline outline-1 outline-sky-950 flex flex-row items-center justify-between px-4">
							<h1>Round Ranking</h1>
							<h1>{roundRanking ?? "-"}/{gameResult.players.length}</h1>
						</div>
					</div>
					<div className="flex justify-center">
						<PrimaryButton size="medium" onClick={handleReturnHome}>Return To Home</PrimaryButton>
					</div>
				</div>
			</div>
		</div>
	)
}
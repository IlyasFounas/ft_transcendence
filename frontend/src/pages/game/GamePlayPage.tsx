import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useLobbySocket } from "../../context/GameSocketContext";
import { useAuth } from "../../context/AuthContext";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import GameBar from "../../components/game/GameBar";
import LoadStreetView from "./StreetView/LoadStreetView";
import Load2DMap from "./2DMap/Load2DMap";
import Preview2DMap from "./2DMap/Preview2DMap";
import TutorialCarousel from "../../components/carousels/TutorialCarousel";

export default function GamePlayPage() {
	const [countdown, setCountDown] = useState(0);
	const [isMapMain, setIsMapMain] = useState(false);
	const { round, roundHistory, lastRoundResult, sendGuess, startGame, leaveGame, resetGame, lobbyInfos, players, isConnected, connect } = useLobbySocket();
	const [guess, setGuess] = useState<{ lat: number; lng: number } | null>(null);
	const navigate = useNavigate();
	const [hasGuessed, setHasGuessed] = useState(false);
	const [showTutorial, setShowTutorial] = useState(false);
	const [gameId, setGameId] = useState<number | null>(null);
	const {authFetch, user} = useAuth();
	const [mode, setMode] = useState<"solo" | "multi" | null>(null);
	const isLeavingRef = useRef(false);

	useEffect(() => {
		setHasGuessed(false);
		setGuess(null);
	}, [round?.round_nb]);

	useEffect(() => {
		if (isLeavingRef.current) return;
		if (!isConnected) connect();
	}, [isConnected, user?.ID]);

	useEffect(() => {
		if (!round) return;

		const endsAt = new Date(round.round_ends_at).getTime();

		const updateCountdown = () => {
			const remaining = Math.max(0, Math.round((endsAt - Date.now()) / 1000));
			setCountDown(remaining);
		};

		updateCountdown();
		const timer = setInterval(updateCountdown, 1000);

		return () => clearInterval(timer);
	}, [round]);

	useEffect(() => {
		authFetch("/api/lobby/current").then((res) => {
			if (res?.status === 200 && res.data?.game_id) {
				setGameId(res.data.game_id);
			}
		});
	}, []);

	useEffect(() => {
		if (!gameId) return;
		authFetch(`/api/game/${gameId}/mode`).then((res) => {
			if (res?.status === 200 && res.data?.mode) {
				setMode(res.data.mode);
			}
		});
	}, [gameId]);

	useEffect(() => {
		if (!lastRoundResult || !round || !mode) return;
		
		const isMultiplayer = mode === "multi";

		if (isMultiplayer) {
			navigate("/results/multiplayers", {
				state: {
					currentRound: round.round_nb,
					totalRounds: lobbyInfos?.round_total ?? 5,
					allResults: lastRoundResult.allResults,
					actualLocation: lastRoundResult.actualLocation,
					nextRoundAt: lastRoundResult.nextRoundAt,
					isLastRound: lastRoundResult.isLastRound,
					fullHistory: lastRoundResult.fullHistory
				}
			});
		} else {
			navigate("/results/solo", {
				state: {
					currentRound: round.round_nb,
					totalRounds: lobbyInfos?.round_total ?? 5,
					rounds: roundHistory,
					actualLocation: lastRoundResult.actualLocation
				}
			});
		}
	}, [lastRoundResult, round, roundHistory, navigate, mode, lobbyInfos?.round_total]);

	useEffect(() => {
		if (lobbyInfos && !round) 
			startGame();
	}, [lobbyInfos, round, startGame]);
	
	const handleSubmit = () => {
		if (!guess || !round) 
			return;
		sendGuess(guess);
		setHasGuessed(true);
	};

	const handleLeave = () => {
		isLeavingRef.current = true;
		leaveGame();
		resetGame();
		navigate("/");
	}

	const handleTuto = () => {
		setShowTutorial(true);
	}

	if (!round) {
		return <div className="h-screen w-full flex items-center justify-center text-white">Loading mission...</div>;
	}

	return (
		<div className="flex flex-col h-screen">
			<div >
				<GameBar
					gameId={gameId ?? 0}
					totalRounds={lobbyInfos?.round_total ?? 5}
					totalTime={lobbyInfos?.timer ?? 120}
					round={round.round_nb}
					countdown={countdown}
					onLeave={handleLeave}
					onTuto={handleTuto}
				/>
			</div>

			<div className="relative flex flex-1 lg:px-4 lg:pb-4">
				{isMapMain && (
					<div className={`
						flex
						w-full h-full
						rounded-lg
						outline outline-1 outline-sky-950
						overflow-hidden
						${hasGuessed ? "pointer-events-none opacity-80" : ""}
					`}>
						<Load2DMap 
							onGuess={(lat, lng) => {
								if (!hasGuessed) {
									setGuess({ lat, lng });
								}
							}} 
							initialGuess={guess} 
						/>
					</div>
				)}
				<div className={
					isMapMain
						? `
							absolute bottom-[1vh] right-[10px]
							lg:bottom-[10px] lg:right-[10px]
							z-1000
							flex flex-col items-center gap-1
							w-[150px] lg:w-[220px]
						`
						: `
							flex
							w-full h-full
							rounded-lg
							outline outline-1 outline-sky-950
							overflow-hidden
						`
				}>
					<div className={
						isMapMain
							? `
								relative shrink-0
								w-[130px] h-[100px]
								lg:w-[190px] lg:h-[160px]
								overflow-hidden
								opacity-75
								cursor-pointer
								transition-transform duration-200
								origin-bottom-right
								ease-out
								hover:scale-[1.1]
								hover:opacity-100
								active:scale-[0.98]
								rounded-[10px]
							`
							: "w-full h-full rounded-lg"
					}>
						<div className={isMapMain ? "absolute inset-0" : "w-full h-full rounded-lg"}>
							<LoadStreetView
								id={String(round.mid)}
								interactive={!isMapMain}
							/>
						</div>
						{isMapMain && (
							<div className="absolute inset-0" onClick={() => setIsMapMain(false)} />
						)}
					</div>
					{isMapMain && (
						<PrimaryButton
							size="small"
							disabled={!guess || hasGuessed}
							onClick={handleSubmit}
						>
							{hasGuessed
								? players.length > 1
									? "Waiting..."
									: "Submitted"
								: "Submit"}
						</PrimaryButton>
					)}
				</div>

				{!isMapMain && (
					<div className="
						absolute bottom-[1vh] right-[10px]
						lg:bottom-[30px] lg:right-[10px]
						z-1000
						flex flex-col items-center gap-1
						w-[150px] lg:w-[220px]
					">
						<div
							className="
								relative shrink-0
								w-[130px] h-[100px]
								lg:w-[190px] lg:h-[160px]
								overflow-hidden
								opacity-75
								cursor-pointer
								transition-transform duration-200
								origin-bottom-right
								ease-out
								hover:scale-[1.1]
								hover:opacity-100
								active:scale-[0.98]
								rounded-[10px]
							"
							onClick={() => setIsMapMain(true)}
						>
							<Preview2DMap />
						</div>
					</div>
				)}
			</div>
			
			{showTutorial && (
				<div className= "fixed inset-0 z-[9999]">
				<TutorialCarousel
				mode="game"
				onClose={() => setShowTutorial(false)}/>
				</div>
			)}
			</div>
	);
}
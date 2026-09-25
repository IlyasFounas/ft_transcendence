import { useEffect, useState } from "react";
import { mockGame } from "../../mocks/game.mock";
import { useNavigate } from "react-router-dom";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import GameBar from "../../components/game/GameBar";
import Countdown from "../../components/game/Countdown";
import LoadStreetView from "./StreetView/LoadStreetView";
import Load2DMap from "./2DMap/Load2DMap";
import Preview2DMap from "./2DMap/Preview2DMap";
import { useLobbySocket } from "../../context/GameSocketContext";

export default function GamePlayPage() {
	const [countdown, setCountDown] = useState(0);
	const [isMapMain, setIsMapMain] = useState(false);
	const { round, roundHistory, lastRoundResult, sendGuess, startGame, lobbyInfos } = useLobbySocket();
	const [guess, setGuess] = useState<{ lat: number; lng: number } | null>(null);
	const navigate = useNavigate();

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
		if (!lastRoundResult || !lobbyInfos || !round) return;
		
		navigate("/soloresults", {
			state: {
				currentRound: round.round_nb,
				totalRounds: lobbyInfos.round_total,
				rounds: roundHistory,
				actualLocation: lastRoundResult.actualLocation
			}
		});
	}, [lastRoundResult, lobbyInfos, round, roundHistory, navigate]);

	useEffect(() => {
		if (lobbyInfos && !round) startGame();
	}, [lobbyInfos, round]);

	const handleSubmit = () => {
		if (!guess || !round) return;
		sendGuess(guess);
	};

	if (!lobbyInfos && !round) return <div className="text-white p-10">Chargement de la partie et du round...</div>;
	if (!lobbyInfos) return <div className="text-white p-10">Attente des données de la partie (WebSocket)...</div>;
	if (!round) return <div className="text-white p-10">Attente du signal de round (WebSocket)...</div>;

	return (
		<div className="flex flex-col h-screen">
			<div className="hidden lg:flex">
				<GameBar
					gameId={0}
					totalRounds={lobbyInfos.round_total}
					totalTime={lobbyInfos.timer}
					round={round.round_nb}
					countdown={countdown}
				/>
			</div>

			<div className="relative flex flex-1 lg:px-4 lg:pb-4">
				{isMapMain && (
					<div className="
						flex
						w-full h-full
						rounded-lg
						outline outline-1 outline-sky-950
						overflow-hidden
					">
						<Load2DMap {...mockGame.location} onGuess={(lat, lng) => setGuess({ lat, lng })} initialGuess={guess} />
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
							: "w-full h-full"
					}>
						<div className={isMapMain ? "absolute inset-0" : "w-full h-full"}>
							<LoadStreetView
								id={String(round.mid)}
								moving={mockGame.streetView.moving}
								zooming={mockGame.streetView.zooming}
								panning={mockGame.streetView.panning}
								interactive={!isMapMain}
							/>
						</div>
						{isMapMain && (
							<div className="absolute inset-0" onClick={() => setIsMapMain(false)} />
						)}
					</div>
					{isMapMain && (
						<PrimaryButton size="small" disabled={!guess} onClick={handleSubmit}>
							Submit
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
							<Preview2DMap {...mockGame.location}/>
						</div>
					</div>
				)}
				<div className="absolute left-1/2 -translate-x-1/2 top-5 lg:hidden z-1000">
					<Countdown countdown={countdown}/>
				</div>
			</div>
		</div>
	);
}
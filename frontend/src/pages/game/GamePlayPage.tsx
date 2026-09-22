import { useEffect, useState, useRef } from "react";
import { mockGame } from "../../mocks/game.mock";
import { useAuth } from "../../context/AuthContext";
import { useParams } from "react-router-dom";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import GameBar from "../../components/game/GameBar";
import Countdown from "../../components/game/Countdown";
import LoadStreetView from "./StreetView/LoadStreetView";
import Load2DMap from "./2DMap/Load2DMap";
import Preview2DMap from "./2DMap/Preview2DMap";

export default function GamePlayPage() {

	const [countdown, setCountDown] = useState(mockGame.totalTime);
	const [isMapMain, setIsMapMain] = useState(false);
	const { gameId } = useParams();
	const { authFetch, user } = useAuth();
	const [game, setGame] = useState<{ ID: number; Timer: number; RoundTotal: number } | null>(null);
	const [round, setRound] = useState<{ round_nb: number; mid: number } | null>(null);
	const [guess, setGuess] = useState<{ lat: number; lng: number } | null>(null);
	const wsRef = useRef<WebSocket | null>(null);

	useEffect(() => {
		if (!user) return;
		const protocol = window.location.protocol === "https:" ? "wss" : "ws";
		const ws = new WebSocket(`${protocol}://${window.location.host}/api/ws/lobby/join?user_id=${user.ID}`);
		wsRef.current = ws;

		ws.onopen = () => {
			ws.send(JSON.stringify({ Action: "start_round", Payload: "" }));
		};

		ws.onmessage = (event) => {
			const message = JSON.parse(event.data);
			if (message.type === "round_started") {
				setRound({ round_nb: message.data.round_nb, mid: message.data.mid });
			}
		};

		return () => ws.close();
	}, [user]);

	useEffect(() => {
		const loadGame = async () => {
			const result = await authFetch(`/api/game/${gameId}`);
			if (result?.status === 200) {
				setGame(result.data.game);
				setCountDown(result.data.game.Timer);
			}
		};
		loadGame();
	}, [gameId]);

	useEffect(() => {
		if (countdown <= 0)
			return;

		const timer = setInterval(() => {
			setCountDown((prev) => prev - 1);
		}, 1000);

		return () => clearInterval(timer);
	}, [countdown]);

	const handleSubmit = () => {
		if (!guess || !round || !wsRef.current) return;

		wsRef.current.send(JSON.stringify({
			Action: "guess",
			Payload: {
				lat: guess.lat,
				long: guess.lng,
				round_nb: round.round_nb,
			},
		}));
	};

	if (!game || !round) return null;

	return (
		<div className="flex flex-col h-screen">
			<div className="hidden lg:flex">
				<GameBar
					gameId={game.ID}
					totalRounds={game.RoundTotal}
					totalTime={game.Timer}
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
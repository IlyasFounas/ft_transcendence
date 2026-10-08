import RoomCodeCard from "../../components/cards/RoomCodeCard";
import LobbyList from "../../components/cards/LobbyList";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import SecondaryButton from "../../components/buttons/SecondaryButton";
import MissionSettingsCard from "../../components/cards/MissionSettingsCard";
import WaitingDots from "../../components/lobby/WaitingDots";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useLobbySocket } from "../../context/GameSocketContext";
import NextRoundTimer from "../../components/results/NextRoundTimer";

export default function LobbyRoomPage() {

	const { players, leaveLobby, lobbyInfos, kickPlayer, isCurrentUserHost, isClosed, startGame, round, gameStartAt, resetGame } = useLobbySocket();
	const navigate = useNavigate();

	useEffect(() => {
  		if (isClosed) navigate("/home");
	}, [isClosed, navigate]);

	useEffect(() => {
		if (round) {
			navigate("/game/multi"); 
		}
	}, [round, navigate]);

	function handleLeave() {
   		leaveLobby();
		resetGame();
		navigate("/home");
  	}

	return (
		<div className="shrink-0 min-h-[calc(100vh-5rem)] w-full bg-cover bg-no-repeat bg-center bg-[url(/backgrounds/earth_background.png)] lg:pr-[27rem]">
		<div className="mx-auto w-full max-w-[1400px] grid grid-cols-1 px-6 py-8 gap-6 lg:grid-cols-[minmax(0,1fr)_394px] lg:items-start lg:gap-8">
			<div className="order-1 lg:col-start-2 lg:row-start-1">
			<RoomCodeCard code={lobbyInfos?.join_code ?? ""} />
			</div>

			<div className="order-2 lg:col-start-2 lg:row-start-2">
			<MissionSettingsCard timeLimit={lobbyInfos?.timer ?? 0} rounds={lobbyInfos?.round_total ?? 0} />
			</div>

			<div className="order-3 lg:col-start-1 lg:row-span-3">
			<LobbyList numberOfPlayers={players.length} listOfPlayers={players} onKick={kickPlayer} isCurrentUserHost={isCurrentUserHost && !gameStartAt}/>
			</div>

			<div className="order-4 flex w-full flex-col items-center gap-4 lg:col-start-2 lg:row-start-3 lg:max-w-none">
			<SecondaryButton 
				variant="danger" 
				className={`w-full ${gameStartAt ? "opacity-50 pointer-events-none" : ""}`}
				onClick={handleLeave}>
				Leave Room
			</SecondaryButton>
			
			{gameStartAt ? (
				<div className="w-full bg-sky-950/40 rounded-[10px] outline outline-1 outline-sky-950 overflow-hidden">
					<NextRoundTimer targetTime={gameStartAt}>Mission starts in</NextRoundTimer>
				</div>
			) : isCurrentUserHost ? (
				<PrimaryButton size="medium" onClick={startGame} disabled={players.length < 2}>
					Start Mission
				</PrimaryButton>
			) : (
				<PrimaryButton 
					size="medium" 
					disabled 
					className="disabled:opacity-100!"
					subtitle="Waiting for the host to start"
				>
					<WaitingDots />
				</PrimaryButton>
			)}
			</div>
		</div>
		</div>
	);
}
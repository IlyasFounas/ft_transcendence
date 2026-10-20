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
		<div className="
			shrink-0 
			min-h-[calc(100vh-5rem)] 
			w-full 
			bg-cover bg-no-repeat bg-center bg-[url(/backgrounds/earth_background.png)] 
			xl:pr-[27rem]">
		<div className="mx-auto w-full max-w-[1400px] grid grid-cols-1 px-6 py-8 gap-6 xl:grid-cols-[minmax(0,1fr)_394px] xl:items-start xl:gap-8">
			<div className="order-1 xl:col-start-2 xl:row-start-1">
			<RoomCodeCard code={lobbyInfos?.join_code ?? ""} />
			</div>

			<div className="order-2 xl:col-start-2 xl:row-start-2">
			<MissionSettingsCard timeLimit={lobbyInfos?.timer ?? 0} rounds={lobbyInfos?.round_total ?? 0} />
			</div>

			<div className="order-3 xl:col-start-1 xl:row-span-3">
			<LobbyList numberOfPlayers={players.length} listOfPlayers={players} onKick={kickPlayer} isCurrentUserHost={isCurrentUserHost && !gameStartAt}/>
			</div>

			<div className="order-4 flex w-full flex-col items-center gap-4 xl:col-start-2 xl:row-start-3 xl:max-w-none">
			<SecondaryButton 
				variant="danger" 
				className={`w-full ${gameStartAt ? "opacity-50 pointer-events-none" : ""}`}
				onClick={handleLeave}>
				Leave Room
			</SecondaryButton>
			
			{gameStartAt ? (
				<NextRoundTimer targetTime={gameStartAt}>Mission starts in</NextRoundTimer>
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
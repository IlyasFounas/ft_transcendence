import RoomCodeCard from "../../components/cards/RoomCodeCard";
import LobbyList from "../../components/cards/LobbyList";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import SecondaryButton from "../../components/buttons/SecondaryButton";
import MissionSettingsCard from "../../components/cards/MissionSettingsCard";
import WaitingDots from "../../components/lobby/WaitingDots";
import { mockGame } from "../../mocks/game.mock";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLobbySocket } from "../../hooks/useLobbySocket";

export default function LobbyPage() {

	const [isReady, setIsReady] = useState(false);
	const { players, leaveLobby } = useLobbySocket();
	const navigate = useNavigate();

	function handleLeave() {
   		leaveLobby();
    	navigate("/home");
  	}

	return (
		<div className="min-h-[calc(100vh-5rem)] w-full bg-[url(/backgrounds/earth_background.png)] lg:pr-[27rem]">
		<div className="mx-auto w-full max-w-[1400px] grid grid-cols-1 px-6 py-8 gap-6 lg:grid-cols-[minmax(0,1fr)_394px] lg:items-start lg:gap-8">
			<div className="order-1 lg:col-start-2 lg:row-start-1">
			<RoomCodeCard code={mockGame.roomcode} />
			</div>

			<div className="order-2 lg:col-start-2 lg:row-start-2">
			<MissionSettingsCard />
			</div>

			<div className="order-3 lg:col-start-1 lg:row-span-3">
			<LobbyList numberOfPlayers={players.length} listOfPlayers={players} />
			</div>

			<div className="order-4 flex w-full flex-col items-center gap-4 lg:col-start-2 lg:row-start-3 lg:max-w-none">
			<SecondaryButton variant="danger" className="w-full" onClick={handleLeave}>Leave Room</SecondaryButton>
			<div>

			</div>
			<PrimaryButton onClick={() => setIsReady(true)} size="medium">
				{isReady ?  <WaitingDots/> : "Ready"}
			</PrimaryButton>
			</div>
		</div>
		</div>
	);
}
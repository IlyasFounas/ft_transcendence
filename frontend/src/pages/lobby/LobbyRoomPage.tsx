import RoomCodeCard from "../../components/cards/RoomCodeCard";
import LobbyList from "../../components/cards/LobbyList";
import PrimaryButton from "../../components/buttons/PrimaryButton";
import SecondaryButton from "../../components/buttons/SecondaryButton";
import MissionSettingsCard from "../../components/cards/MissionSettingsCard";
import WaitingDots from "../../components/lobby/WaitingDots";
import { mockLobbyPlayers } from "../../mocks/lobbyPlayers.mock";
import { mockGame } from "../../mocks/game.mock";
import { useState } from "react";

export default function LobbyPage() {

	const [isReady, setIsReady] = useState(false);

  return (
    <div className="min-h-screen w-full bg-[url(/backgrounds/earth_background.png)]">
      <div className="container mx-auto grid grid-cols-1 px-6 py-8 gap-6 lg:grid-cols-[minmax(0,1fr)_394px] lg:items-start lg:gap-8">
        <div className="order-1 lg:col-start-2 lg:row-start-1">
          <RoomCodeCard code={mockGame.roomcode} />
        </div>

        <div className="order-2 lg:col-start-2 lg:row-start-2">
          <MissionSettingsCard />
        </div>

        <div className="order-3 lg:col-start-1 lg:row-span-3">
          <LobbyList numberOfPlayers={mockLobbyPlayers.length} listOfPlayers={mockLobbyPlayers} />
        </div>

        <div className="order-4 flex w-full flex-col items-center gap-4 lg:col-start-2 lg:row-start-3 lg:max-w-none">
          <SecondaryButton className="w-full">Leave Room</SecondaryButton>
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
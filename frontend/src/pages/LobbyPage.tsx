import RoomCodeCard from "../components/cards/RoomCodeCard";
import LobbyTestPage from "../components/cards/LobbyTestPage";
import PrimaryButton from "../components/buttons/PrimaryButton";
import SecondaryButton from "../components/buttons/SecondaryButton";
import MissionSettingsCard from "../components/cards/MissionSettingsCard";

export default function LobbyPage()
{
	return (
		<div className="
		min-h-screen w-full 
		bg-[url(/backgrounds/earth_background.png)]
		">
      		<div className="
				container mx-auto 
				grid grid-cols-1
				px-6 py-8 
				gap-6
				
				
				lg:grid-cols-[minmax(0,1fr)_394px]
				lg:items-start
				lg:gap-8
				">
				
				<div className="order-1 lg:col-start-2 lg:row-start-1 ">
        			<RoomCodeCard code="ATL-5463" />
				</div>

				<div className="order-2 lg:col-start-2 lg:row-start-2">
				<MissionSettingsCard />
				</div>

				<div className="order-3 lg:col-start-1 lg:row-span-3">
        		<LobbyTestPage />
				</div>

				<div className="
					order-4 
					flex w-full
					flex-col items-center gap-4
					
					lg:col-start-2
					lg:row-start-3
					lg:max-w-none
				">
					<SecondaryButton className="
						w-full ">
						Leave Room
					</SecondaryButton>

					<PrimaryButton size="medium">
						Ready
					</PrimaryButton>
				</div>
			</div>
    </div>
	)
}
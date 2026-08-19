import GameBar from "../components/game/GameBar";
import LoadStreetView from "./game/StreetView/LoadStreetView";

export default function GamePlayPage()
{
	return (
		<div className="flex flex-col h-screen">
			<GameBar/>
			<div className="
				flex w-full flex-1 border
				rounded-lg ">
				<LoadStreetView id="271688394642360" moving={false} zooming={true} panning={true}></LoadStreetView>
			</div >
		</div>
	)
}
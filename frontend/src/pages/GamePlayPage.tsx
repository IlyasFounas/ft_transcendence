import {useState}  from "react";

import PrimaryButton from "../components/buttons/PrimaryButton";
import GameBar from "../components/game/GameBar";
import LoadStreetView from "./game/StreetView/LoadStreetView";
import Preview2DMap from "./game/2DMap/Preview2DMap";
import CountdownCard from "../components/game/Countdown";
import Load2DMap from "./game/2DMap/Load2DMap";

export default function GamePlayPage()
{
	
	const totalTime = 60;
	const [countDown, setCountDown] = useState(totalTime);
	
	const [isMapMain, setIsMapMain] = useState(false);

return (
	<div className="flex flex-col h-screen">

		<div className="hidden lg:flex">
			<GameBar />
		</div>

		<div className="relative flex flex-1 lg:px-4 lg:pb-4">

			<div className="
				flex
				w-full h-full
				outline outline-1 outline-sky-950 rounded-lg
				overflow-hidden
			">
				{isMapMain ? (
					<Load2DMap
						latitude={45.782429}
						longitude={4.745860}
					/>
				) : (
					<LoadStreetView
						id="271688394642360"
						moving={true}
						zooming={true}
						panning={true}
					/>
				)}
			</div>

			<div className="
				absolute left-1/2 -translate-x-1/2 top-5
				lg:hidden z-1000 
			">
				<CountdownCard
					countDown={countDown}
					setCountDown={setCountDown}
				/>
			</div>

			<div className="
				absolute
				bottom-[1vh] right-[10px]
				lg:bottom-[10px] lg:right-[10px]
				flex flex-col
				items-center gap-1
				w-[150px]
				lg:w-[220px]
				z-1000
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
					onClick={() => setIsMapMain(prev => !prev)}
				>
					{isMapMain ? (
						<LoadStreetView
							id="271688394642360"
							moving={false}
							zooming={false}
							panning={false}
						/>
					) : (
						<Preview2DMap
							latitude={45.782429}
							longitude={4.745860}
						/>
					)}
				</div>

				<PrimaryButton size="small">
					Submit
				</PrimaryButton>

			</div>

		</div>
	</div>
);
}
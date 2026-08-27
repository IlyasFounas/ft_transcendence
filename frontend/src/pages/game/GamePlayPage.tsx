import { useEffect, useState } from "react";

import PrimaryButton from "../../components/buttons/PrimaryButton";
import GameBar from "../../components/game/GameBar";
import Countdown from "../../components/game/Countdown";
import LoadStreetView from "./StreetView/LoadStreetView";
import Load2DMap from "./2DMap/Load2DMap";

import { mockGame } from "../../mocks/game.mock";
import Preview2DMap from "./2DMap/Preview2DMap";

export default function GamePlayPage() {

    const [countdown, setCountDown] = useState(mockGame.totalTime);
    const [isMapMain, setIsMapMain] = useState(false);

	useEffect(() => {
    if (countdown <= 0)
        return;

    const timer = setInterval(() => {
        setCountDown((prev) => prev - 1);
    }, 1000);

    return () => clearInterval(timer);
	}, [countdown]);

	return (
		<div className="flex flex-col h-screen">

			<div className="hidden lg:flex">
				<GameBar
					gameId={mockGame.id}
					round={mockGame.currentRound}
					totalRounds={mockGame.totalRounds}
					countdown={countdown}
					totalTime={mockGame.totalTime}
				/>
			</div>

			<div className="relative flex flex-1 lg:px-4 lg:pb-4">

				<div className="
					flex
					w-full h-full

					rounded-lg
					outline outline-1 outline-sky-950
					overflow-hidden
				">
					{isMapMain ? (
						<Load2DMap {...mockGame.location}/>
					) : (
						<LoadStreetView
							id={mockGame.streetView.id}
							moving={mockGame.streetView.moving}
							zooming={mockGame.streetView.zooming}
							panning={mockGame.streetView.panning}
						/>
					)}
				</div>

				<div className="absolute left-1/2 -translate-x-1/2 top-5 lg:hidden z-1000">
					<Countdown
						countdown={countdown}/>
				</div>

				<div className="
					absolute bottom-[1vh] right-[10px]
					lg:bottom-[10px] lg:right-[10px]
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
						onClick={() => setIsMapMain(prev => !prev)}
					>
						{isMapMain ? (
							<LoadStreetView
								id={mockGame.streetView.id}
								moving={false}
								zooming={false}
								panning={false}
							/>
						) : (
							<Preview2DMap {...mockGame.location}/>
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
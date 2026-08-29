import PrimaryButton from "../../components/buttons/PrimaryButton";
import ResultInfoRow from "../../components/results/ResultInfoRow";
import TotalScoreCard from "../../components/results/TotalScoreCard";
import Preview2DMap from "./2DMap/Preview2DMap";
import GameResultList from "../../components/results/GameResultList";
import { mockGameResult } from "../../mocks/gameResult.mock";

export default function SoloResultPage()
{

	const gameResult = mockGameResult;
	
	  const currentPlayerResult = gameResult.players.find(
		(p) => p.player.username === gameResult.currentUsername
	  );

	return (
			<div className="
				grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px]
				gap-4
				p-4">

			<div className="order-1 lg:col-start-1 lg:row-start-1">
				<ResultInfoRow
					title="Mission Debrief"
					infos="Round 3 of 5"
					variant="primary"
				/>
			</div>

			<div className="
				order-2
				lg:col-start-1 lg:row-start-2
				w-full h-60 lg:h-[650px]
				outline outline-1 outline-sky-950
				rounded-[10px]
				overflow-hidden
			">
				<Preview2DMap latitude={42} longitude={12} />
			</div>

			<div className="
				flex flex-row order-3 gap-2
				lg:col-start-1 lg:row-start-3
			">
				<ResultInfoRow
					title="Total Score"
					infos="4578"
					variant="secondary"
					suffix="pts"
				/>
				<ResultInfoRow
					title="Distance"
					infos="4225"
					variant="secondary"
					suffix="km"
				/>
			</div>
			<div className="hidden lg:flex lg:col-start-2 lg:row-start-1">
				<ResultInfoRow title="Your results"variant="primary"/>
			</div>
			<div className="
				order-4
			   	lg:col-start-2 lg:row-start-2 lg:row-span-3
				flex flex-col
				items-center
				gap-4
			">
				<GameResultList rounds={currentPlayerResult?.rounds ?? []} />
				<TotalScoreCard points={456} distance={9874} />
				<PrimaryButton size="medium">
					Next Round
				</PrimaryButton>
			</div>

			</div>
	)
}
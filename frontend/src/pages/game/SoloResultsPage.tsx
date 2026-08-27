import PrimaryButton from "../../components/buttons/PrimaryButton";
import GameResultList from "../../components/results/GameResultList";
import PanelInfosCard from "../../components/results/PanelInfosCard";
import ScoreResumeCard from "../../components/results/ScoreSummaryCard";
import Load2DMap from "./2DMap/Load2DMap";

export default function SoloResultPage()
{
	return (
			<div className="
				grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px]
				gap-4
				p-4">

			<div className="order-1 lg:col-start-1 lg:row-start-1">
				<PanelInfosCard
					title="Mission Debrief"
					infos="Round 3 of 5"
					style="primary"
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
				<Load2DMap latitude={45.782429} longitude={4.745860} />
			</div>

			<div className="
				flex flex-row order-3 gap-2
				lg:col-start-1 lg:row-start-3
			">
				<PanelInfosCard
					title="Total Score"
					infos="4578"
					style="secondary"
					suffix="pts"
				/>
				<PanelInfosCard
					title="Distance"
					infos="4225"
					style="secondary"
					suffix="km"
				/>
			</div>
			<div className="hidden lg:flex lg:col-start-2 lg:row-start-1">
				<PanelInfosCard title="Your results" style="primary"/>
			</div>
			<div className="
				order-4
			   	lg:col-start-2 lg:row-start-2 lg:row-span-3
				flex flex-col
				items-center
				gap-4
			">
						<GameResultList rounds={5} /><GameResultList rounds={currentPlayerResult?.rounds ?? []} />
				<ScoreResumeCard points={456} distance={9874} />
				<PrimaryButton size="medium">
					Next Round
				</PrimaryButton>
			</div>

			</div>
	)
}
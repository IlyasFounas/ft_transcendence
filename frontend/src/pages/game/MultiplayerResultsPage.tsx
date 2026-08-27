import GameResultList from "../../components/results/GameResultList";
import NextRoundTimer from "../../components/results/NextRoundTimer";
import PanelInfosCard from "../../components/results/PanelInfosCard";
import RankResultsList from "../../components/results/RankResultsList";
import ScoreSummaryCard from "../../components/results/ScoreSummaryCard";
import Load2DMap from "./2DMap/Load2DMap";
import { mockGameResult } from "../../mocks/gameResult.mock";

export default function MultiplayerResultsPage() {
  const gameResult = mockGameResult;

  const currentPlayerResult = gameResult.players.find(
    (p) => p.player.username === gameResult.currentPlayerUsername
  );
  const totalScore =
    currentPlayerResult?.rounds.reduce((sum, r) => sum + r.points, 0) ?? 0;

  return (
    <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_394px] gap-4 p-4">
      <div className="order-1 lg:col-start-1 lg:row-start-1">
        <PanelInfosCard
          title="Mission Debrief"
          infos={`Round ${gameResult.currentRound} of ${gameResult.totalRounds}`}
          style="primary"
        />
      </div>

      <div className="order-2 lg:col-start-1 lg:row-start-2 w-full h-60 lg:h-[380px] outline outline-1 outline-sky-950 rounded-[10px] overflow-hidden">
        <Load2DMap latitude={45.782429} longitude={4.745860} />
      </div>

      <div className="flex flex-row order-3 gap-2 lg:hidden">
        <PanelInfosCard
          title="Total Score"
          infos={String(totalScore)}
          style="secondary"
          suffix="pts"
        />
        <PanelInfosCard
          title="Distance"
          infos={String(gameResult.distance)}
          style="secondary"
          suffix="km"
        />
      </div>

      <div className="hidden lg:flex lg:col-start-2 lg:row-start-1">
        <PanelInfosCard title="Your results" style="primary" />
      </div>

      <div className="order-4">
        <RankResultsList
          players={gameResult.players}
          currentRound={gameResult.currentRound}
        />
      </div>

      <div className="order-5 lg:col-start-2 lg:row-start-2 lg:row-span-3 flex flex-col items-center gap-4">
		<GameResultList rounds={currentPlayerResult?.rounds ?? []} />
        <ScoreSummaryCard points={totalScore} distance={gameResult.distance} />
        <NextRoundTimer/>
      </div>
    </div>
  );
}
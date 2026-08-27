import SummaryResultCard from "./GameSummaryCard";
import { GRID_COLS, ROUND_VISIBILITY } from "./GameSummaryCard";

import type {PlayerResult}  from "../../types/game.ts";

interface RankResultListProps {
  players: PlayerResult[];
  currentRound: number;
}

export default function RankResultList({ players, currentRound }: RankResultListProps) {
  const visibleCount = Math.min(currentRound, ROUND_VISIBILITY.length);
  const displayRounds = Array.from(
    { length: visibleCount },
    (_, i) => currentRound - visibleCount + i + 1
  );
  const offset = ROUND_VISIBILITY.length - displayRounds.length;

  const sortedPlayers = [...players]
    .sort((a, b) => {
      const totalA = a.rounds.reduce((sum, r) => sum + r.points, 0);
      const totalB = b.rounds.reduce((sum, r) => sum + r.points, 0);
      return totalB - totalA;
    })
    .map((result, index) => ({ ...result, rank: index + 1 }));

  return (
    <div className="@container flex flex-col rounded-[12px] border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden w-full">
      <div className="h-[230px] lg:h-[350px] overflow-y-auto [scrollbar-gutter:stable]">
        <div className={`sticky top-0 z-10 h-[40px] grid ${GRID_COLS} items-center px-4
          uppercase text-heading-small text-text-primary/70 bg-surface-base`}>
          <h1 className="text-center">Rank</h1>
          <h1 className="px-8">Agent</h1>
          {displayRounds.map((round, index) => (
            <h1
              key={round}
              className={`text-center justify-center ${ROUND_VISIBILITY[offset + index]}`}
            >
              Round {round}
            </h1>
          ))}
          <h1 className="text-center">Total</h1>
        </div>

        {sortedPlayers.map((result) => (
          <SummaryResultCard
            key={result.player.username}
            rank={result.rank}
            player={result.player}
            rounds={result.rounds}
            displayRounds={displayRounds}
          />
        ))}
      </div>
    </div>
  );
}
import RankResultRow from "./RankResultRow";
import { GRID_COLS, ROUND_VISIBILITY } from "./rankResults.constants.ts";
import { getSortedPlayersWithRank } from "../../utils/ranking.ts";
import type {PlayerResult}  from "../../types/game.ts";

interface RankResultListProps {
	players: PlayerResult[];
	currentRound: number;
	currentUsername: string;
	onHoverPlayer: (username: string | null) => void;
}

export default function RankResultList({ players, currentRound, currentUsername, onHoverPlayer}: RankResultListProps) {
  const visibleCount = Math.min(currentRound, ROUND_VISIBILITY.length);
  const displayRounds = Array.from(
	{ length: visibleCount },
	(_, i) => currentRound - visibleCount + i + 1
  );
  	const offset = ROUND_VISIBILITY.length - displayRounds.length;

	const sortedPlayers = getSortedPlayersWithRank(players);
  return (
	<div className="@container flex flex-col rounded-[12px] border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden w-full">
	  <div className="h-[230px] lg:h-[600px] overflow-y-auto [scrollbar-gutter:stable]">
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
		  <RankResultRow
			key={result.player.username}
			rank={result.rank}
			player={result.player}
			rounds={result.rounds}
			displayRounds={displayRounds}
			currentUsername={currentUsername}
			onMouseEnter={() => onHoverPlayer(result.player.username)}
			onMouseLeave={() => onHoverPlayer(null)}
		  />
		))}
	  </div>
	</div>
  );
}
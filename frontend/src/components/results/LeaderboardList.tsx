import LeaderboardRow from './LeaderboardRow.tsx'
import { getSortedPlayersWithRank } from "../../utils/ranking.ts";
import type {PlayerResult}  from "../../types/game.ts";
import LeaderboardRowPlayer from './LeaderboardRowPlayer.tsx';

type RankedPlayerResult = ReturnType<typeof getSortedPlayersWithRank>[number];

interface LeaderBoardListProps {
  players: RankedPlayerResult[];
  currentPlayerResult?: RankedPlayerResult;
  currentUsername: string;
  onHoverPlayer?: (username: string | null) => void;
}

export default function LeaderBoardList({ players, currentPlayerResult, currentUsername, onHoverPlayer}: LeaderBoardListProps) {

  return (
	<div className="flex flex-col gap-4 lg:w-[620px]">
		<div className="flex flex-col rounded-[12px] border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden w-full">
		<div className=" ">
			<div className={`z-10 h-[50px] grid grid-cols-[30%_40%_30%] items-center
			uppercase text-heading-small text-text-primary/70 bg-surface-base`}>
			<h1 className="text-center px-8">Rank</h1>
			<h1 className="px-8">Agent</h1>
			<h1 className="text-center"> Score</h1>
			</div>
			<div className="lg:h-[500px] overflow-y-auto [scrollbar-gutter:stable]">

			{players.map((result) => (
				<LeaderboardRow
					key={result.player.username}
					rank={result.rank}
					player={result.player}
					rounds={result.rounds}
					currentUsername={currentUsername}
					onMouseEnter={() => onHoverPlayer?.(result.player.username)}
					onMouseLeave={() => onHoverPlayer?.(null)}
				/>
			))}
			</div>
		</div>
		</div>
		{currentPlayerResult && (
			<LeaderboardRowPlayer
				rank={currentPlayerResult.rank}
				player={currentPlayerResult.player}
				score={currentPlayerResult.rounds}
			/>
		)}
	</div>
  );
}
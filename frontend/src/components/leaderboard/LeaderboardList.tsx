import LeaderboardRow from './LeaderboardRow.tsx'
import LeaderboardRowPlayer from './LeaderboardRowPlayer.tsx';
import type { LeaderboardEntry } from "../../services/leaderboard.ts";
import LeaderboardEmptyRow from './LeaderboardEmptyRow.tsx';

interface LeaderBoardListProps {
  players: LeaderboardEntry[];
  currentPlayerResult?: LeaderboardEntry;
  currentUsername: string;
  onHoverPlayer?: (username: string | null) => void;
  totalRank: number;
}

export default function LeaderBoardList({ players, currentPlayerResult, currentUsername, onHoverPlayer, totalRank}: LeaderBoardListProps) {

	const rowCount = Math.max(0, totalRank - 3);

  return (
	<div className="flex flex-col gap-4 w-full">
	  <div className="flex flex-col rounded-[12px] border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden w-full">
		<div className=" ">
		  <div className="z-10 h-[50px] grid grid-cols-[30%_40%_30%] items-center
			uppercase text-heading-small text-text-primary/70 bg-surface-base bg-gradient-to-t from-bg-surface-base via-black">
			<h1 className="text-center px-8">Rank</h1>
			<h1 className="px-8">Agent</h1>
			<h1 className="text-center"> Score</h1>
		  </div>
		  <div className=" h-[250px] md:h-[455px] overflow-y-auto custom-scrollbar bg-surface-base">

			{Array.from({ length: rowCount }).map((_, index) => {
				const player = players[index];

				return player ? (
				<LeaderboardRow
					key={player.username}
					position={index + 4}
					username={player.username}
					points={player.total}
					rank={player.rank}
					avatar={player.avatar}
					currentUsername={currentUsername}
					onMouseEnter={() => onHoverPlayer?.(player.username)}
					onMouseLeave={() => onHoverPlayer?.(null)}
			  />
				) : (
					<LeaderboardEmptyRow
						key={`empty-row-${index + 4}`}
						position={index + 4}
					/>
				)
			})} 
		  </div>
		</div>
	  </div>
	  {currentPlayerResult && (
		<LeaderboardRowPlayer
		  position={currentPlayerResult.position}
		  username={currentPlayerResult.username}
		  points={currentPlayerResult.total}
		  avatar={currentPlayerResult.avatar}
		  rank={currentPlayerResult.rank}
		/>
	  )}
	</div>
  );
}
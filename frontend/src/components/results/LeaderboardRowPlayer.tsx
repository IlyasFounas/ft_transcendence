import type { Player, RoundResult } from "../../types/game";
import Avatar from "../primitives/Avatar";

interface LeaderboardRowPlayer {
	player: Player;
	rank: number;
	score: RoundResult[];
}

export default function LeaderboardRowPlayer({
	player,
	rank,
	score,
} : LeaderboardRowPlayer)
{

	return (
		<div
			className="
				uppercase
				text-primary
				w-full h-[60px] lg:h-[80px]
				flex
				grid grid-cols-[30%_40%_30%]
				items-center
				bg-surface-base
				shadow-[0px_0px_25px_0px_rgba(53,223,255,0.50)]

				outline outline-2 rounded-xl
				outline-offset-[-2px] 
				outline-cyan-400/60"
			>
			  <div className="text-center text-heading-medium">{rank}</div>
		
			  <div className="flex items-center gap-[12px] @[640px]:px-[30px] min-w-0">
				<Avatar avatar={player.avatar} size={45} />
				<div className="flex flex-col min-w-0">
				  <span className={`text-heading-small truncate`}>
					{player.username}
				  </span>
				  <span className="hidden @[860px]:block text-body-primary truncate">
					{player.role}
				  </span>
				</div>
			  </div>
		
			  <div className="flex items-center justify-center gap-[10px]">
				<span className="text-heading-small">
				  {score.reduce((total, r) => total + r.points, 0)}
				</span>
				<span className="text-body-primary">pts</span>
			  </div>
		</div>
	)
}
import Avatar from "../primitives/Avatar";
import type { Player, RoundResult } from "../../types/game";

type Props = {
	rank: number;
	player: Player;
	rounds: RoundResult[];
	currentUsername: string;
	onMouseEnter?: () => void;
  	onMouseLeave?: () => void;
};

export default function LeaderboardRow({
	rank = 2,
	player,
	rounds,
	currentUsername,
	onMouseEnter, 
	onMouseLeave
}: Props) {

  const isCurrentPlayer = player.username === currentUsername;

  const rowColor = isCurrentPlayer ? "text-cyan-400" : "text-(--color-text-secondary)";
  const avatarBorder = isCurrentPlayer ? "#35DFFF" : "var(--color-text-secondary)";

  return (
	<div
		onMouseEnter={onMouseEnter}
		onMouseLeave={onMouseLeave}
		className={`
			w-full 
			h-[65px] 
			outline outline-1 outline-offset-[-1px] 
			bg-surface-base
			items-center 
			uppercase
			grid grid-cols-[30%_40%_30%]
			${rowColor}
			hover:bg-cyan-400/8 transition-colors
			${
			isCurrentPlayer
				? "outline-cyan-400/60 bg-gradient-to-r from-cyan-400/15 via-cyan-400/5 to-transparent shadow-[0_0_20px_rgba(53,223,255,0.25)]"
				: "outline-cyan-800/30"
			}`}
	>
	  <div className="text-center text-heading-medium">{rank}</div>
	  <div className="flex items-center gap-[12px]">
		<Avatar avatar={player.avatar} borderColor={avatarBorder} size={45} />
		<div className="flex flex-col min-w-0">
		  <span className={`text-heading-small truncate ${rowColor}`}>
			{player.username}
		  </span>
		  <span className="text-body-primary truncate">
			{player.role}
		  </span>
		</div>
	  </div>

	  <div className="flex items-center justify-center gap-[10px]">
		<span className="text-heading-small">
		  {rounds.reduce((total, r) => total + r.points, 0)}
		</span>
		<span className="text-body-primary">pts</span>
	  </div>
	</div>
  );
}
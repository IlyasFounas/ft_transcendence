import Avatar from "../primitives/Avatar";
import { getAgentColorByIndex } from "../../utils/agentColor";

type Props = {
  position: number;
  username: string;
  rank: string,
  role?: string;
  avatar?: string;
  points: number;
  currentUsername: string;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
};

export default function LeaderboardRow({
  position,
  username,
  rank,
  avatar,
  points,
  currentUsername,
  onMouseEnter,
  onMouseLeave
}: Props) {

  const isCurrentPlayer = username === currentUsername;
  const agentColor = isCurrentPlayer ? "#35DFFF" : getAgentColorByIndex(position - 1);

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
			hover:bg-cyan-400/8 transition-colors
			${
			isCurrentPlayer
			? "outline-cyan-400/60 bg-gradient-to-r from-cyan-400/15 via-cyan-400/5 to-transparent shadow-[0_0_20px_rgba(53,223,255,0.25)]"
			: "outline-cyan-800/30"
		}`}
	  style={{ color: agentColor }}
	>
	  <div className="text-center text-heading-medium">{position}</div>
	  <div className="flex items-center gap-[12px]">
		<Avatar avatar={avatar} borderColor={agentColor} size={45} />
		<div className="flex flex-col min-w-0">
		  <span className={`text-heading-small truncate text-white`}>
			{username}
		  </span>
		  <span className="text-body-primary truncate">
			{rank}
		  </span>
		</div>
	  </div>

	  <div className="flex items-center justify-center gap-[10px]">
		<span className="text-heading-small">
		  {points}
		</span>
		<span className="text-body-primary">pts</span>
	  </div>
	</div>
  );
}
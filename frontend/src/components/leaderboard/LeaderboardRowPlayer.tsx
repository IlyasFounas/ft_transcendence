import Avatar from "../primitives/Avatar";

interface LeaderboardRowPlayerProps {
  username: string;
  rank?: string;
  avatar?: string;
  position: number;
  points: number;
}

export default function LeaderboardRowPlayer({
	username,
	rank,
	avatar,
	position,
	points,
} : LeaderboardRowPlayerProps)
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
	  <div className="text-center text-heading-medium">{position}</div>
	  <div className="flex items-center gap-[12px] @[640px]:px-[30px] min-w-0">
		<div className="shrink-0">
			<Avatar avatar={avatar} size={45} />
		</div>
		<div className="flex flex-col min-w-0">
		  <span className="text-heading-small truncate text-text-primary">
			{username}
		  </span>
		  <span className="hidden md:inline text-body-primary truncate">
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
	)
}
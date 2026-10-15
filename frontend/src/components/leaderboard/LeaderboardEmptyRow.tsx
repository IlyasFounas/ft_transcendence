import Avatar from "../primitives/Avatar";

interface EmptyLeaderboardRowProps {
    position: number;
}

export default function LeaderboardEmptyRow({
    position
}: EmptyLeaderboardRowProps)
{
     return (
        <div
          className="
            w-full 
			h-[65px] 
			outline outline-1 outline-offset-[-1px] 
			bg-surface-base
			items-center 
			uppercase
			grid grid-cols-[30%_40%_30%]
			hover:bg-cyan-400/8 transition-colors
            outline-cyan-800/30"
        >
          <div className="text-center text-heading-medium">{position}</div>
          <div className="flex items-center gap-[12px] @[640px]:px-[30px] min-w-0">
            <Avatar size={45} borderColor="#07192E"/>
            <div className="flex flex-col min-w-0">
              <span className={`text-heading-small truncate`}>
                -
              </span>
            </div>
          </div>
          <div className="flex items-center justify-center gap-[10px]">
            <span className="text-heading-small">
              -
            </span>
            <span className="text-body-primary">pts</span>
          </div>
        </div>
     )
}
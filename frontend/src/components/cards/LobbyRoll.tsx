import Avatar from "../primitives/Avatar.tsx"
import StatusBadge from '../primitives/StatusBadge.tsx'
import type {StatusBadgeProps}  from "../primitives/StatusBadge";


interface LobbyRollProps {
	username?: string;
	role?: string;
	avatar: string;
	status?: StatusBadgeProps["status"];
}

export default function LobbyRoll({
		username,
		role,
		status,
		avatar
}: LobbyRollProps) {

	
	return (
		<div className="w-full min-h-[65px] 
		outline outline-1 outline-offset-[-1px] outline-cyan-800/30
		flex items-center justify-between px-[16px] uppercase">
  			<div className="flex items-center gap-[12px]">
			<Avatar size={45} avatar={avatar}/>
			<div className="flex flex-col px-[10px]">
	  			<span className="text-heading-small text-text-primary">
	   				{username}
	 			</span>
	  			<span className="text-body-primary text-primary">
					{role}
	  			</span>
			</div>
  		</div>
 		<StatusBadge status={status ?? "waiting"} />
		</div>
	);
}


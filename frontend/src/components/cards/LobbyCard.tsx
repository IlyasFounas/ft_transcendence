import Avatar from "../primitives/Avatar.tsx"
import StatusBadge from '../primitives/StatusBadge.tsx'

interface LobbyCardProps {
	username?: string;
	role?: string;
}

export default function LobbyCard({
		username = "Agent Tommy",
		role = "Field Analyst"
}: LobbyCardProps) {

	
	return (
		<div className="w-[480px] h-[65px] 
		outline outline-1 outline-offset-[-1px] outline-cyan-800/30
		flex items-center justify-between px-[16px] uppercase">
  			<div className="flex items-center gap-[12px]">
			<Avatar size="45"/>
			<div className="flex flex-col px-[10px]">
	  			<span className="text-heading-small text-text-primary">
	   				{username}
	 			</span>
	  			<span className="text-body-primary text-primary">
					{role}
	  			</span>
			</div>
  		</div>
 		<StatusBadge status="ready" />
		</div>
	);
}


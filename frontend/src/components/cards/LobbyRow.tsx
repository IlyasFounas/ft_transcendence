import Avatar from "../primitives/Avatar.tsx"
import StatusBadge from '../primitives/StatusBadge.tsx'
import type {StatusBadgeProps}  from "../primitives/StatusBadge.tsx";


interface LobbyRowProps {
	username?: string;
	role?: string;
	avatar: string;
	status?: StatusBadgeProps["status"];
	currentUsername: string; 
}

export default function LobbyRow({
		username,
		role,
		status,
		avatar,
		currentUsername,
}: LobbyRowProps) {

		const isCurrentPlayer = username === currentUsername;
	 	const rowColor = isCurrentPlayer ? "text-cyan-400" : "text-(--color-text-secondary)";
  		const avatarBorder = isCurrentPlayer ? "#35DFFF" : "var(--color-text-secondary)";

	return (
		<div className={`
		w-full min-h-[65px] 
		bg-surface-base
		outline outline-1 outline-offset-[-1px] outline-cyan-800/30
		flex items-center 
		justify-between 
		px-[16px] 
		uppercase
		${rowColor}
		${
		  isCurrentPlayer
			? "outline-cyan-400/60 bg-gradient-to-r from-cyan-400/15 via-cyan-400/5 to-transparent shadow-[0_0_20px_rgba(53,223,255,0.25)]"
			: "outline-cyan-800/30"
		}`}>
  			<div className="flex items-center gap-[12px]">
			<Avatar size={45} avatar={avatar} borderColor={avatarBorder}/>
			<div className="flex flex-col px-[10px]">
	  			<span className={`text-heading-small ${rowColor}`}>
	   				{username}
	 			</span>
	  			<span className={`text-body-primary ${rowColor}`}>
					{role}
	  			</span>
			</div>
  		</div>
 		<StatusBadge status={status ?? "waiting"} />
		</div>
	);
}


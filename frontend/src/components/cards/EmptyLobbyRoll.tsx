import Avatar from "../primitives/Avatar.tsx"
import StatusBadge from '../primitives/StatusBadge.tsx'
import type {StatusBadgeProps}  from "../primitives/StatusBadge";

export default function EmptyLobbyRoll() {

	
	return (
		<div className="
			w-full min-h-[65px] 
			bg-surface-base
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center 
			justify-between 
			px-[16px] 
			uppercase
		">
  			<div className="flex items-center gap-[12px]">
			<Avatar size={45} borderColor="#07192E"/>
  		</div>
		</div>
	);
}


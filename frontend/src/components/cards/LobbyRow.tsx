import Avatar from "../primitives/Avatar.tsx"
import { X } from "lucide-react"

interface LobbyRowProps {
	username?: string;
	avatar: string;
	currentUsername: string;
	isHost?: boolean;
	isCurrentUserHost?: boolean;
	onKick?: () => void;
}

export default function LobbyRow({
		username,
		avatar,
		currentUsername,
		isHost,
		isCurrentUserHost,
		onKick,
}: LobbyRowProps) {

		const isCurrentPlayer = username === currentUsername;
	 	const rowColor = isCurrentPlayer ? "text-cyan-400" : "text-(--color-text-secondary)";
  		const avatarBorder = isCurrentPlayer ? "#35DFFF" : "#6E8AA5";
		

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
  			<div className="flex items-center gap-[12px] min-w-0">
			<div className="shrink-0">
				<Avatar size={45} avatar={avatar} borderColor={avatarBorder}/>
			</div>
			<div className="flex flex-col px-[10px] min-w-0">
				<span className={`truncate text-heading-small ${rowColor}`}>
					{username}{isHost && " (host)"}
				</span>
			</div>
		</div>
 		{isCurrentUserHost && !isHost && 
		<button onClick={onKick}>
			<X className="
				cursor-pointer
				transition-transform
				duration-200
				ease-out
				hover:scale-[1.2]
				active:scale-[0.8]
				"/>
		</button>}
		</div>
	);
}

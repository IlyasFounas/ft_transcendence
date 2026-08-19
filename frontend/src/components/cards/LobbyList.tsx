import LobbyRoll from './LobbyRoll'
import EmptyLobbyRoll from './EmptyLobbyRoll'

import type { StatusBadgeProps } from "../primitives/StatusBadge";

interface Player {
  id: string;
  username: string;
  role: string;
  status: StatusBadgeProps["status"];
  avatar: string;
}

interface LobbyListProps{
	className?: string;
	numberOfPlayers: number;
	listOfPlayers: Player[];
}

export default function LobbyList({className, listOfPlayers}: LobbyListProps)
{

	const playersReady = listOfPlayers.filter((player) => player.status == "ready").length;
	const numberOfPlayers = listOfPlayers.length;

	return (
		<div 
		className={`
			rounded-[12px] 
			border 
			border-border-default/30
            shadow-[0_0_10px_rgba(40,215,255,0.5)] 
			overflow-hidden 
			${className}`}>
      			<div className="
					flex 
					items-center 
					justify-between 
					border-b 
					border-border-default/30 
                  	bg-surface-title px-4 py-3">
        				<span className="text-heading-small text-text-primary/70 uppercase">
        					Agents
        				</span>
						<span>
							{playersReady}/{numberOfPlayers}
						</span>
      			</div>
			<div className="flex flex-col max-h-[350px] lg:max-h-full overflow-y-auto">
  				{Array.from({ length: 10 }).map((_, index) => {
				const player = listOfPlayers[index];

				return player ? (
					<LobbyRoll
						key={player.id}
						username={player.username}
						role={player.role}
						status={player.status}
						avatar={player.avatar}
					/>
				) : (
					<EmptyLobbyRoll key={`empty-${index}`} />
				);
			})}
			</div>
   		</div>
	)
}

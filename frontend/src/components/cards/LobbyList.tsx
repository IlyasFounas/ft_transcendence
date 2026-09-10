import LobbyRow from './LobbyRow'
import EmptyLobbyRow from './EmptyLobbyRow'

import type { StatusBadgeProps } from "../primitives/StatusBadge";

interface Player {
  id: string;
  username: string;
  role: string;
  status: StatusBadgeProps["status"];
  avatar: string;
}

interface LobbyListProps{
	numberOfPlayers: number;
	listOfPlayers: Player[];
}

export default function LobbyList({listOfPlayers}: LobbyListProps)
{
	const playersReady = listOfPlayers.filter((player) => player.status == "ready").length;
	const numberOfPlayers = listOfPlayers.length;

	return (
		<div 
		className="
			rounded-[12px] 
			border 
			border-border-default/30
            shadow-[0_0_10px_rgba(40,215,255,0.5)] 
			overflow-hidden ">
      			<div className="
					flex 
					bg-surface-base
					items-center 
					justify-between 
					border
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
					<LobbyRow
						key={player.id}
						index={index}
						username={player.username}
						role={player.role}
						status={player.status}
						avatar={player.avatar}
						currentUsername="Agent Tommy"
					/>
				) : (
					<EmptyLobbyRow key={`empty-${index}`} />
				);
			})}
			</div>
   		</div>
	)
}

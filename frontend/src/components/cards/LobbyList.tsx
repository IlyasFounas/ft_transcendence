import LobbyRow from './LobbyRow'
import EmptyLobbyRow from './EmptyLobbyRow'
import { useAuth } from '../../context/AuthContext';

interface Player {
  id: number;
  username: string;
  avatar: string;
  isHost?: boolean;
}

interface LobbyListProps{
	numberOfPlayers: number;
	listOfPlayers: Player[];
	onKick: (playerId: number) => void;
	isCurrentUserHost?: boolean;
}

export default function LobbyList({listOfPlayers, onKick, isCurrentUserHost}: LobbyListProps)
{
	const numberOfPlayers = listOfPlayers.length;
	const { user } = useAuth();

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
					bg-gradient-to-t from-black to-bg-surface-base 
                  	bg-surface-title 
					px-4 py-3
					uppercase">
        				<span className="text-heading-small text-text-primary/70">
        					Agents
        				</span>
						<span>
  						{numberOfPlayers} {numberOfPlayers > 1 ? "agents" : "agent"}
						</span>
      			</div>
			<div className="flex flex-col max-h-[350px] lg:max-h-full overflow-y-auto">
  				{Array.from({ length: 10 }).map((_, index) => {
				const player = listOfPlayers[index];

				return player ? (
					<LobbyRow
						key={player.id}
						username={player.username}
						avatar={player.avatar}
						currentUsername={user?.Username ?? ""}
						isHost={player.isHost}
						isCurrentUserHost={isCurrentUserHost}
						onKick={() => onKick(player.id)}
					/>
				) : (
					<EmptyLobbyRow key={`empty-${index}`} />
				);
			})}
			</div>
   		</div>
	)
}
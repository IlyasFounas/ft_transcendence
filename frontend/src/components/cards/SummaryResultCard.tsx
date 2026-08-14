import Avatar from "../primitives/Avatar";

interface Player {
	username: string;
	role: string;
	avatar: string;
}

interface RoundResult {
	round: number;
	points: number;
}

interface GameSummaryCardProps {
	rank: number;
	rounds: RoundResult[]
	player: Player;
}

export default function GameSummaryCard({
	rank = 2,
	rounds,
	player,
}: GameSummaryCardProps) {


	return (
		<div className="
			w-[1015px] h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center px-[16px] justify-between
			uppercase text-primary
		">
			<div className="w-[60px] text-center text-heading-medium px-[5px]">
				{rank}
			</div>
			<div className="flex items-center gap-[12px] px-[50px]">
						<Avatar size={45}/>
						<div className="flex flex-col px-[10px]">
							<span className="text-heading-small text-text-primary">
								{player.username}
							</span>
							<span className="text-body-primary text-primary">
								{player.role}
							</span>
						</div>
			</div>
			{rounds.map((round) => (
  			<div key={round.round} className="flex items-center gap-[10px] px-[14px]">
   				 <div className="text-heading-small">
      				{round.points}
    		</div>

    		<div className="text-body-primary">
      			pts
    		</div>
  			</div>
			))}
		</div>
	);
}
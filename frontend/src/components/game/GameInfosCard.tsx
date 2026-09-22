interface GameInfosCardProps
{
	totalRounds: number;
	round: number;
	id: number;
}

export default function GameInfosCard({round, totalRounds, id}: GameInfosCardProps) 
{
	return (
		<div className="
			w-full h-[75px]
			flex flex-col justify-center
			rounded-lg
			px-2
			shadow-[inset_0px_0px_20px_0px_rgba(0,175,255,0.05)] outline outline-1 outline-sky-950
			uppercase
 			">
			<h1 className="
				text-heading-large
				text-primary">Operation #{id}</h1>
			<h2 className="
				text-heading-small
				text-text-secondary">Round <span className="text-primary">
						 {round}
					</span> of {totalRounds}</h2>
		</div>
	)
}
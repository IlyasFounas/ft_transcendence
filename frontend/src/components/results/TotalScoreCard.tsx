interface TotalScoreCardProps{
	points: number;
	distance: number;
}

export default function TotalScoreCard({points, distance} : TotalScoreCardProps)
{
	return (
		<div className="
			uppercase
			text-primary
			w-full h-[60px] lg:h-[80px]
			flex
			flex-row
			justify-around
			items-center
			shadow-[0px_0px_25px_0px_rgba(53,223,255,0.50)]

			outline outline-2 rounded-xl
			outline-offset-[-2px] 
			outline-cyan-400/60">
			<h1 className="text-heading-medium">Total</h1>
			<h1 className="text-heading-small">{points} pts</h1>
			<h1 className="text-heading-small">{distance} km</h1>
		</div>
	)
}
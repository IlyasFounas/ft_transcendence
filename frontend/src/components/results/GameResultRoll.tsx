interface GameResultCardProps {
	round: number;
	points: number;
	distance: number;
}

export default function GameResultList({
	round,
	points,
	distance,
}: GameResultCardProps) {

	const stats = [
		{ value: points, label: "pts" },
		{ value: distance, label: "km" },
	];

	return (
		<div className="
			w-full h-[60px] lg:h-[80px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex justify-around items-center
			uppercase text-primary
		">
			<div className="w-[60px] text-center text-heading-medium">
				{round}
			</div>

			{stats.map((stat) => (
				<div 
					key={stat.label}
					className="flex items-center gap-[10px] p-[14px]"
				>
					<div className="text-heading-small">
						{stat.value}
					</div>

					<div className="text-body-primary">
						{stat.label}
					</div>
				</div>
			))}
		</div>
	);
}
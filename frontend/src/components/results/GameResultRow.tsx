import { formatDistance } from "../../utils/formatDistance";

interface GameResultRowProps {
	round: number;
	points: number | null;
	distance: number | null;
	currentRound: number;
}

export default function GameResultRow({
	round,
	points,
	distance,
	currentRound,
}: GameResultRowProps) {

	const formattedDistance = distance !== null ? formatDistance(distance) : null;

	const stats = [
		{ value: points, label: "pts" },
		{ value: formattedDistance?.value ?? null, label: formattedDistance?.unit ?? "km" },
	];

	const isCurrentRound = currentRound === round;
	const rowColor = isCurrentRound ? "text-cyan-400" : "text-(--color-text-secondary)";

	return (
		<div className={`
			w-full h-[60px] lg:h-[100px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex justify-around items-center
			uppercase ${rowColor}
			${
				isCurrentRound
					? "outline-cyan-400/60 bg-gradient-to-r from-cyan-400/15 via-cyan-400/5 to-transparent shadow-[0_0_20px_rgba(53,223,255,0.25)]"
					: "outline-cyan-800/30"
			}
		`}>
			<div className="w-[60px] text-center text-heading-medium">
				{round}
			</div>

			{stats.map((stat) => (
				<div
					key={stat.label}
					className="flex items-center gap-[10px] p-[14px]"
				>
					<div className="text-heading-small">
						{stat.value !== null ? stat.value : "-"}
					</div>

					<div className="text-body-primary">
						{stat.label}
					</div>
				</div>
			))}
		</div>
	);
}
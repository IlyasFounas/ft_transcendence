function polarToCartesian(
  centerX: number,
  centerY: number,
  radius: number,
  angleInDegrees: number
) {
  const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180.0;

  return {
    x: centerX + radius * Math.cos(angleInRadians),
    y: centerY + radius * Math.sin(angleInRadians),
  };
}

function describeArc(
  centerX: number,
  centerY: number,
  radius: number,
  startAngle: number,
  endAngle: number
): string {
  const start = polarToCartesian(centerX, centerY, radius, endAngle);
  const end = polarToCartesian(centerX, centerY, radius, startAngle);

  const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";

  return [
    "M", start.x, start.y,
    "A", radius, radius, 0, largeArcFlag, 0, end.x, end.y,
  ].join(" ");
}

const GAP = 3;
const SIZE = 100;
const CENTER = SIZE / 2;
const RADIUS = 30;
const STROKE_WIDTH = 5;
const TOTAL_SEGMENTS = 10;

export function CountdownRing({countDown} : {countDown : number}) {
  
	const segmentAngle = 360 / TOTAL_SEGMENTS;
	

  return (
	<div className="relative" style={{ width: SIZE, height: SIZE }}>
		<svg width={SIZE} height={SIZE} viewBox={`0 0 ${SIZE} ${SIZE}`}>
			{Array.from({ length: TOTAL_SEGMENTS }).map((_, index) => {
			const startAngle = index * segmentAngle + GAP / 2;
			const endAngle = (index + 1) * segmentAngle - GAP / 2;
			const isVisible = index < countDown;

			return (
				<path
				key={index}
				d={describeArc(CENTER, CENTER, RADIUS, startAngle, endAngle)}
				stroke="currentColor"
				strokeWidth={STROKE_WIDTH}
				strokeLinecap="butt"
				fill="none"
				className={isVisible ? "text-primary" : "text-surface-elevated"}
				/>
			);
			})}
		</svg>

		<div className="absolute inset-0 flex items-center justify-center">
			<span className="text-heading-large text-primary">
			{countDown}
			</span>
		</div>
		</div>
  );
}
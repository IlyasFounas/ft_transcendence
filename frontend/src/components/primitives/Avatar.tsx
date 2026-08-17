import { useId } from "react";

type AvatarSvgProps = {
  size?: number;
  borderColor?: string;
  avatar?: string;
  className?: string;
};

const HEXAGON_PATH =
  "M129.5 43.793V110.206L70 143.427L10.5 110.206L10.5 43.793L70 10.5723L129.5 43.793Z";

function hexToRgbNormalized(hex: string): [number, number, number] {
  const r = parseInt(hex.slice(1, 3), 16) / 255;
  const g = parseInt(hex.slice(3, 5), 16) / 255;
  const b = parseInt(hex.slice(5, 7), 16) / 255;
  return [r, g, b];
}

export default function Avatar({
  size = 60,
  borderColor = "#35DFFF",
  avatar,
  className
	}: AvatarSvgProps) {
	const width = size;
	const height = size * 1.1;

	const uid = useId().replace(/:/g, "");
	const filterId = `hex-avatar-glow-${uid}`;
	const clipId = `hex-avatar-clip-${uid}`;

	const [r, g, b] = hexToRgbNormalized(borderColor);

	return (
		<svg width={width} height={height} viewBox="0 0 140 154" fill="none" xmlns="http://www.w3.org/2000/svg">
			<defs>
				<clipPath id={clipId}>
				<path d={HEXAGON_PATH} />
				</clipPath>

				<filter id={filterId} x="0" y="0" width="140" height="154" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
				<feFlood floodOpacity="0" result="BackgroundImageFix" />
				<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha" />
				<feOffset />
				<feGaussianBlur stdDeviation="5" />
				<feComposite in2="hardAlpha" operator="out" />
				<feColorMatrix type="matrix" values={`0 0 0 0 ${r} 0 0 0 0 ${g} 0 0 0 0 ${b} 0 0 0 1 0`} />
				<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_872_967" />
				<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_872_967" result="shape" />
				</filter>
			</defs>

			<g filter={`url(#${filterId})`}>
				<path d={HEXAGON_PATH} stroke={borderColor} strokeWidth="5" />
			</g>

			<image
				href={avatar}
				x="10.5"
				y="10.5"
				width="119"
				height="133"
				preserveAspectRatio="xMidYMid slice"
				clipPath={`url(#${clipId})`}
			/>

			<path d={HEXAGON_PATH} stroke={borderColor} strokeWidth="5" shapeRendering="crispEdges" />
		</svg>
	);
}
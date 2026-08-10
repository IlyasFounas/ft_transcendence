import { useId } from "react";

type AvatarSvgProps = {
	size?: number;
	borderColor?: string;
	avatar?: string;
};

export default function Avatar({
	size = 60,
	borderColor = "#35DFFF",
	avatar = "avatars/avatar1.png",
}: AvatarSvgProps)
{
	const width = size;
	const height = size * 1.1;

	const uid = useId().replace(/:/g, "");
  	const filterId = `hex-avatar-glow-${uid}`;
  	const clipId = `hex-avatar-clip-${uid}`;

	return (
		<svg width={width} height={height} viewBox="0 0 140 154" fill="none" xmlns="http://www.w3.org/2000/svg">
			<g filter="url(#filter0_d_872_967)">
			<path d="M129.5 43.793V110.206L70 143.427L10.5 110.206L10.5 43.793L70 10.5723L129.5 43.793Z" stroke={borderColor} strokeWidth="5"/>
			</g>
			<defs>
				<clipPath id={clipId}>
		  			<path d="M129.5 43.793V110.206L70 143.427L10.5 110.206L10.5 43.793L70 10.5723L129.5 43.793Z" />
				</clipPath>
			<filter id="filter0_d_872_967" x="0" y="0" width="140" height="154" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB">
			<feFlood flood-opacity="0" result="BackgroundImageFix"/>
			<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
			<feOffset/>
			<feGaussianBlur stdDeviation="5"/>
			<feComposite in2="hardAlpha" operator="out"/>
			<feColorMatrix type="matrix" values="0 0 0 0 0.207843 0 0 0 0 0.87451 0 0 0 0 1 0 0 0 1 0"/>
			<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_872_967"/>
			<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_872_967" result="shape"/>
			</filter>
			</defs>

			<image
				href={avatar}
				x="10.5"
				y="10.5"
				width="119"
				height="133"
				preserveAspectRatio="xMidYMid slice"
				clipPath={`url(#${clipId})`}
			/>

			 <g filter={`url(#${filterId})`}>
				<path d="M129.5 43.793V110.206L70 143.427L10.5 110.206L10.5 43.793L70 10.5723L129.5 43.793Z" stroke={borderColor} shapeRendering="crispEdges" />
	  		</g>
		</svg>
	)
}
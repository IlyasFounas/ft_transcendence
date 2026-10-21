import { useId } from "react";

type AvatarHexagonProps = {
	avatar?: string;
	className?: string;
};

export function BronzeAvatarHexagon({
	avatar = "avatars/avatar1.png",
	className,
	}: AvatarHexagonProps)
{
	const uid = useId().replace(/:/g, "");
	const filterId = `hex-avatar-glow-${uid}`;
	const clipId = `hex-avatar-clip-${uid}`;


	return (
		<svg className={className} width="100" height="110" viewBox="0 0 100 110" fill="none" xmlns="http://www.w3.org/2000/svg">
			<g filter="url(#filter0_d_412_477)">
			<path d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z" stroke="url(#paint0_linear_412_477)" strokeWidth="4" shapeRendering="crispEdges"/>
			</g>
			<defs>
			<filter id="filter0_d_412_477" x="0" y="0" width="99.282" height="110" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
			<feFlood floodOpacity="0" result="BackgroundImageFix"/>
			<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
			<feOffset/>
			<feGaussianBlur stdDeviation="7.5"/>
			<feComposite in2="hardAlpha" operator="out"/>
			<feColorMatrix type="matrix" values="0 0 0 0 0.921569 0 0 0 0 0.580392 0 0 0 0 0.298039 0 0 0 1 0"/>
			<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_412_477"/>
			<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_412_477" result="shape"/>
			</filter>
			<linearGradient id="paint0_linear_412_477" x1="82.488" y1="65.3786" x2="33.6421" y2="34.4683" gradientUnits="userSpaceOnUse">
			<stop offset="0.0576923" stopColor="#5E321E"/>
			<stop offset="0.641451" stopColor="#CE561E"/>
			<stop offset="1" stopColor="#EB944C"/>
			</linearGradient>

			<clipPath id={clipId}>
        		<path d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z" />
        	</clipPath>
			</defs>

		{avatar && (
			<image
					href={avatar}
					x="17"
					y="17"
					width="65.28"
					height="75.4"
					preserveAspectRatio="xMidYMid slice"
					clipPath={`url(#${clipId})`}
			/>
			)}

			<g filter={`url(#${filterId})`}>
				<path
				d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z"
				stroke="url(#paint0_linear_412_477)"
				strokeWidth="4"
				shapeRendering="crispEdges"
				/>
      		</g>
		</svg>
	)
}

export function SilverAvatarHexagon({
	avatar = "avatars/avatar1.png",
	className,
	}: AvatarHexagonProps)
{
	const uid = useId().replace(/:/g, "");
	const filterId = `hex-avatar-glow-${uid}`;
	const clipId = `hex-avatar-clip-${uid}`;

	return (
		<svg className={className} width="100" height="110" viewBox="0 0 100 110" fill="none" xmlns="http://www.w3.org/2000/svg">
			<g filter="url(#filter0_d_412_447)">
			<path d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z" stroke="url(#paint0_linear_412_447)" strokeWidth="4" shapeRendering="crispEdges"/>
			</g>
			<defs>
			<filter id="filter0_d_412_447" x="0" y="0" width="99.282" height="110" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
			<feFlood floodOpacity="0" result="BackgroundImageFix"/>
			<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
			<feOffset/>
			<feGaussianBlur stdDeviation="7.5"/>
			<feComposite in2="hardAlpha" operator="out"/>
			<feColorMatrix type="matrix" values="0 0 0 0 0.767792 0 0 0 0 0.80272 0 0 0 0 0.809295 0 0 0 1 0"/>
			<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_412_447"/>
			<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_412_447" result="shape"/>
			</filter>
			<linearGradient id="paint0_linear_412_447" x1="72.6057" y1="96.8652" x2="50.3545" y2="44.6134" gradientUnits="userSpaceOnUse">
			<stop stopColor="#475967"/>
			<stop offset="1" stopColor="#C9DBE9"/>
			</linearGradient>

			<clipPath id={clipId}>
        		<path d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z" />
        	</clipPath>
			</defs>

		{avatar && (
			<image
					href={avatar}
					x="17"
					y="17"
					width="65.28"
					height="75.4"
					preserveAspectRatio="xMidYMid slice"
					clipPath={`url(#${clipId})`}
			/>
		)}

			<g filter={`url(#${filterId})`}>
				<path
				d="M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z"
				stroke="url(#paint0_linear_412_447)"
				strokeWidth="4"
				shapeRendering="crispEdges"
				/>
      		</g>
		</svg>
	)
}

export function GoldAvatarHexagon({
  avatar = "avatars/avatar1.png",
  className,
}: AvatarHexagonProps) {
  const uid = useId().replace(/:/g, "");
  const filterId = `hex-avatar-glow-${uid}`;
  const clipId = `hex-avatar-clip-${uid}`;

  return (
    <svg className={className} width="104" height="115" viewBox="0 0 104 115" fill="none" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <filter id={filterId} x="0" y="0" width="103.612" height="115" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
          <feFlood floodOpacity="0" result="BackgroundImageFix" />
          <feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha" />
          <feOffset />
          <feGaussianBlur stdDeviation="7.5" />
          <feComposite in2="hardAlpha" operator="out" />
          <feColorMatrix type="matrix" values="0 0 0 0 1 0 0 0 0 0.905155 0 0 0 0 0.38782 0 0 0 1 0" />
          <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_412_462" />
          <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_412_462" result="shape" />
        </filter>

        <clipPath id={clipId}>
          <path d="M86.6118 37.4043V77.5947L51.8061 97.6895L17.0004 77.5947V37.4043L51.8061 17.3096L86.6118 37.4043Z" />
        </clipPath>
      </defs>

      {avatar && (
		<image
				href={avatar}
				x="17"
				y="17"
				width="65.28"
				height="75.4"
				preserveAspectRatio="xMidYMid slice"
				clipPath={`url(#${clipId})`}
		/>
		)}

      <g filter={`url(#${filterId})`}>
        <path
          d="M86.6118 37.4043V77.5947L51.8061 97.6895L17.0004 77.5947V37.4043L51.8061 17.3096L86.6118 37.4043Z"
          stroke="#DDB74E"
          strokeWidth="4"
          shapeRendering="crispEdges"
        />
      </g>
    </svg>
  );
}
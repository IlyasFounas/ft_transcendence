import { useId } from "react";
import { AVATAR_FRAME_VARIANTS } from "./AvatarFrameConfig";

export function AvatarFrame({
  variant, size = 100, imageSrc,
}: { variant: "gold" | "silver" | "bronze"; size?: number; imageSrc: string }) {
  const config = AVATAR_FRAME_VARIANTS[variant];
  const uid = useId().replace(/:/g, "");
  const filterId = `avatar-glow-${uid}`;
  const clipId = `avatar-clip-${uid}`;
  const gradientId = `avatar-gradient-${uid}`;

  const displayWidth = size;
  const displayHeight = size * (config.viewBoxHeight / config.viewBoxWidth);
  const strokeValue = typeof config.stroke === "string" ? config.stroke : `url(#${gradientId})`;

	const PAD = 1.15;
	const { x: bx, y: by, width: bw, height: bh } = config.hexBounds;

	const imageWidth = bw * PAD;
	const imageHeight = bh * PAD;
	const imageX = bx - (imageWidth - bw) / 2;
	const imageY = by - (imageHeight - bh) / 2;	

  return (
    <svg width={displayWidth} height={displayHeight} viewBox={`0 0 ${config.viewBoxWidth} ${config.viewBoxHeight}`} fill="none" xmlns="http://www.w3.org/2000/svg">
      <g filter={`url(#${filterId})`}>
        <path d={config.hexPath} stroke={strokeValue} strokeWidth={config.strokeWidth} />
      </g>
      <defs>
        <clipPath id={clipId}>
          <path d={config.hexPath} />
        </clipPath>
        <filter id={filterId} x={0} y={0} width={config.viewBoxWidth} height={config.viewBoxHeight} filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
          <feFlood floodOpacity="0" result="bg" />
          <feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha" />
          <feOffset />
          <feGaussianBlur stdDeviation={config.blur} />
          <feComposite in2="hardAlpha" operator="out" />
          <feColorMatrix type="matrix" values={`0 0 0 0 ${config.glowColor[0]} 0 0 0 0 ${config.glowColor[1]} 0 0 0 0 ${config.glowColor[2]} 0 0 0 ${config.glowOpacity} 0`} />
          <feBlend mode="normal" in2="bg" result="glow" />
          <feBlend mode="normal" in="SourceGraphic" in2="glow" result="shape" />
        </filter>
        {typeof config.stroke !== "string" && (
          <linearGradient id={gradientId} {...config.stroke.coords} gradientUnits="userSpaceOnUse">
            {config.stroke.stops.map((s, i) => <stop key={i} offset={s.offset} stopColor={s.color} />)}
          </linearGradient>
        )}
      </defs>
      <image
        href={imageSrc}
        x={imageX}
        y={imageY}
        width={imageWidth}
        height={imageHeight}
        preserveAspectRatio="xMidYMid slice"
        clipPath={`url(#${clipId})`}
      />
    </svg>
  );
}
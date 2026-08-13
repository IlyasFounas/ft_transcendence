import { useId, Fragment } from "react";
import { PODIUM_VARIANTS } from "./PodiumShapeConfig";

function GlowFilter({
  id, x, y, width, height, layers,
}: { id: string; x: number; y: number; width: number; height: number; layers: readonly { blur: number; color: [number, number, number]; opacity: number }[] }) {
  return (
    <filter id={id} x={x} y={y} width={width} height={height} filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
      <feFlood floodOpacity="0" result="bg" />
      {layers.map((layer, i) => {
        const prev = i === 0 ? "bg" : `${id}-glow-${i - 1}`;
        const result = `${id}-glow-${i}`;
        return (
          <Fragment key={i}>
            <feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha" />
            <feOffset />
            <feGaussianBlur stdDeviation={layer.blur} />
            <feComposite in2="hardAlpha" operator="out" />
            <feColorMatrix type="matrix" values={`0 0 0 0 ${layer.color[0]} 0 0 0 0 ${layer.color[1]} 0 0 0 0 ${layer.color[2]} 0 0 0 ${layer.opacity} 0`} />
            <feBlend mode="normal" in2={prev} result={result} />
          </Fragment>
        );
      })}
      <feBlend mode="normal" in="SourceGraphic" in2={`${id}-glow-${layers.length - 1}`} result="shape" />
    </filter>
  );
}

export function PodiumShapeSvg({ variant, className }: { variant: "gold" | "silver" | "bronze"; className?: string }) {
  const config = PODIUM_VARIANTS[variant];
  const uid = useId().replace(/:/g, "");
  const outerFilterId = `podium-outer-filter-${uid}`;
  const innerFilterId = `podium-inner-filter-${uid}`;
  const outerGradientId = `podium-outer-gradient-${uid}`;
  const innerGradientId = `podium-inner-gradient-${uid}`;

  return (
    <svg className={className} width={config.width} height={config.height} viewBox={`0 0 ${config.width} ${config.height}`} fill="none" xmlns="http://www.w3.org/2000/svg">
      <g opacity="0.8" filter={`url(#${outerFilterId})`}>
        <path d={config.outer.fillPath} fill="#020A14" />
        <path d={config.outer.strokePath} stroke={`url(#${outerGradientId})`} strokeWidth={2} />
      </g>
      <g opacity="0.8" filter={`url(#${innerFilterId})`}>
        <path d={config.inner.fillPath} fill="#020A14" />
        <path d={config.inner.strokePath} stroke={`url(#${innerGradientId})`} />
      </g>
      <defs>
        <GlowFilter id={outerFilterId} x={0} y={0} width={config.width} height={config.height} layers={config.outer.glow} />
        <GlowFilter id={innerFilterId} x={5} y={5} width={config.width - 10} height={config.height - 11} layers={config.inner.glow} />
        <linearGradient id={outerGradientId} {...config.outer.gradientCoords} gradientUnits="userSpaceOnUse">
          {config.outer.gradientStops.map((s, i) => <stop key={i} offset={s.offset} stopColor={s.color} />)}
        </linearGradient>
        <linearGradient id={innerGradientId} {...config.inner.gradientCoords} gradientUnits="userSpaceOnUse">
          {config.inner.gradientStops.map((s, i) => <stop key={i} offset={s.offset} stopColor={s.color} />)}
        </linearGradient>
      </defs>
    </svg>
  );
}
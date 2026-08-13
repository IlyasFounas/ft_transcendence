type AvatarStroke = string | { stops: { offset?: number; color: string }[]; coords: { x1: number; y1: number; x2: number; y2: number } };

type AvatarFrameConfig = {
  viewBoxWidth: number; viewBoxHeight: number;
  hexPath: string;
  strokeWidth: number;
  blur: number;
  glowColor: [number, number, number];
  glowOpacity: number;
  stroke: AvatarStroke;
  hexBounds: {x: number, y: number, width: number, height: number};
};

export const AVATAR_FRAME_VARIANTS: Record<"gold" | "silver" | "bronze", AvatarFrameConfig> = {
  gold: {
    viewBoxWidth: 104, viewBoxHeight: 115,
    hexPath: "M86.6117 37.4043V77.5947L51.806 97.6895L17.0004 77.5947V37.4043L51.806 17.3096L86.6117 37.4043Z",
    strokeWidth: 6, blur: 7.5, glowColor: [1, 0.905155, 0.38782], glowOpacity: 1,
    stroke: "#DDB74E",
	hexBounds: { x: 17.0004, y: 17.3096, width: 69.6113, height: 80.3799 },
  },
  silver: {
    viewBoxWidth: 100, viewBoxHeight: 110,
    hexPath: "M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z",
    strokeWidth: 4, blur: 7.5, glowColor: [0.767792, 0.80272, 0.809295], glowOpacity: 1,
    stroke: { stops: [{ color: "#475967" }, { offset: 1, color: "#C9DBE9" }], coords: { x1: 72.6057, y1: 96.8652, x2: 50.3545, y2: 44.6134 } },
	hexBounds: { x: 17.0004, y: 17.3096, width: 65.2812, height: 75.3799 }
  },
  bronze: {
    viewBoxWidth: 100, viewBoxHeight: 110,
  	hexPath: "M82.2816 36.1543V73.8447L49.641 92.6895L17.0004 73.8447V36.1543L49.641 17.3096L82.2816 36.1543Z", // 👈 remplacé, identique au silver
    strokeWidth: 4, blur: 7.5, glowColor: [0.921569, 0.580392, 0.298039], glowOpacity: 1,
    stroke: { stops: [{ offset: 0.0576923, color: "#5E321E" }, { offset: 0.641451, color: "#CE561E" }, { offset: 1, color: "#EB944C" }], coords: { x1: 97.488, y1: 80.3786, x2: 48.6421, y2: 49.4683 }},
	hexBounds: { x: 17.0004, y: 17.3096, width: 65.2812, height: 75.3799 }
  },
};

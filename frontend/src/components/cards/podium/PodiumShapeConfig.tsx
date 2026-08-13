// podiumShapeConfig.ts
type GlowLayer = { blur: number; color: [number, number, number]; opacity: number };
type GradientStop = { offset?: number; color: string };
type GradientCoords = { x1: number; y1: number; x2: number; y2: number };

type PodiumLayer = {
  fillPath: string;
  strokePath: string;
  glow: [GlowLayer, GlowLayer, GlowLayer];
  gradientStops: GradientStop[];
  gradientCoords: GradientCoords;
};

type PodiumVariantConfig = { width: number; height: number; outer: PodiumLayer; inner: PodiumLayer };

export const PODIUM_VARIANTS: Record<"gold" | "silver" | "bronze", PodiumVariantConfig> = {
  gold: {
    width: 312, height: 422,
    outer: {
      fillPath: "M45 44.7501L54 36.0001H258L267.377 45.1165L276 53.5001V366.313L267 376.735L259 386H53L44.6529 376.333L36 366.313V53.5001L45 44.7501Z",
      strokePath: "M257.594 37L266.68 45.834L275 53.9219V365.94L266.243 376.082L258.543 385H53.457L45.4102 375.68L37 365.94V53.9219L45.6973 45.4668L54.4062 37H257.594Z",
      glow: [
        { blur: 2, color: [1, 0.964706, 0.74902], opacity: 0.9 },
        { blur: 8, color: [1, 0.847059, 0.290196], opacity: 0.45 },
        { blur: 18, color: [1, 0.721569, 0], opacity: 0.2 },
      ],
      gradientStops: [
        { offset: 0.000671374, color: "#D6B58F" },
        { offset: 0.47529, color: "#FFD45A" },
        { offset: 0.999771, color: "#DACBB0" },
      ],
      gradientCoords: { x1: 115, y1: 55.3334, x2: 245.047, y2: 398.107 },
    },
    inner: {
      fillPath: "M49.625 49.4766L58.25 41H253.75L262.736 49.8316L271 57.9531V360.99L262.375 371.087L254.708 380.062H57.2917L49.2924 370.698L41 360.99V57.9531L49.625 49.4766Z",
      strokePath: "M253.546 41.5L262.386 50.1885L270.5 58.1621V360.806L261.995 370.763L254.478 379.562H57.5225L49.6729 370.373L41.5 360.805V58.1621L49.9756 49.833L58.4541 41.5H253.546Z",
      glow: [
        { blur: 2, color: [1, 0.964706, 0.74902], opacity: 0.9 },
        { blur: 8, color: [1, 0.847059, 0.290196], opacity: 0.45 },
        { blur: 18, color: [1, 0.721569, 0], opacity: 0.2 },
      ],
      gradientStops: [
        { offset: 0, color: "#E5B76D" },
        { offset: 1, color: "#C5B58E" },
      ],
      gradientCoords: { x1: 271, y1: 210.531, x2: 40.9936, y2: 213.244 },
    },
  },
  silver: {
    width: 287, height: 372,
    outer: {
      fillPath: "M44.0625 43.5L52.125 36H234.875L243.275 43.8141L251 51V319.125L242.938 328.059L235.771 336H51.2292L43.7516 327.714L36 319.125V51L44.0625 43.5Z",
      strokePath: "M234.482 37L242.594 44.5459L250 51.4346V318.74L242.195 327.389L235.326 335H51.6738L44.4941 327.044L37 318.74V51.4346L44.7432 44.2324L52.5176 37H234.482Z",
      glow: [
        { blur: 2, color: [0.631373, 0.768627, 0.843137], opacity: 0.25 },
        { blur: 8, color: [0.447639, 0.667206, 0.841346], opacity: 0.1 },
        { blur: 18, color: [0.278431, 0.34902, 0.403922], opacity: 0.25 },
      ],
      gradientStops: [
        { offset: 0.000671374, color: "white" },
        { offset: 0.47529, color: "#345C7B" },
        { offset: 1, color: "#667187" },
      ],
      gradientCoords: { x1: 106.771, y1: 36, x2: 242.106, y2: 331.956 },
    },
    inner: {
      fillPath: "M48.7266 48.2656L56.4531 41H231.589L239.639 48.5699L247.042 55.5312V315.277L239.315 323.932L232.447 331.625H55.5946L48.4286 323.598L41 315.277V55.5312L48.7266 48.2656Z",
      strokePath: "M231.391 41.5L239.296 48.9346L246.542 55.7471V315.086L238.942 323.599L232.224 331.125H55.8184L48.8018 323.266L41.5 315.086V55.7471L49.0693 48.6299L56.6514 41.5H231.391Z",
      glow: [
        { blur: 2, color: [0.633937, 0.805762, 0.899038], opacity: 1 },
        { blur: 8, color: [0.207843, 0.87451, 1], opacity: 1 },
        { blur: 18, color: [0.501961, 0.627451, 0.72549], opacity: 1 },
      ],
      gradientStops: [
        { offset: 0, color: "#BFCDDA" },
        { offset: 0.983218, color: "#3B4E6E" },
      ],
      gradientCoords: { x1: 247.221, y1: 150.313, x2: 39.0114, y2: 278.836 },
    },
  },
  bronze: {
    width: 287, height: 372,
    outer: {
      fillPath: "M44.0625 43.5L52.125 36H234.875L243.275 43.8141L251 51V319.125L242.938 328.059L235.771 336H51.2292L43.7516 327.714L36 319.125V51L44.0625 43.5Z",
      strokePath: "M234.482 37L242.594 44.5459L250 51.4346V318.74L242.195 327.389L235.326 335H51.6738L44.4941 327.044L37 318.74V51.4346L44.7432 44.2324L52.5176 37H234.482Z",
      glow: [
        { blur: 2, color: [0.164706, 0.0509804, 0], opacity: 0.25 },
        { blur: 8, color: [0.929412, 0.843137, 0.639216], opacity: 0.2 },
        { blur: 18, color: [0.937255, 0.831373, 0.576471], opacity: 0.3 },
      ],
      gradientStops: [
        { offset: 0.000671374, color: "#EBAE5A" },
        { offset: 0.47529, color: "#9F6348" },
        { offset: 1, color: "#87340F" },
      ],
      gradientCoords: { x1: 106.771, y1: 36, x2: 242.106, y2: 331.956 },
    },
    inner: {
      fillPath: "M48.7266 48.2656L56.4531 41H231.589L239.639 48.5699L247.042 55.5312V315.277L239.315 323.932L232.447 331.625H55.5946L48.4286 323.598L41 315.277V55.5312L48.7266 48.2656Z",
      strokePath: "M231.391 41.5L239.296 48.9346L246.542 55.7471V315.086L238.942 323.599L232.224 331.125H55.8184L48.8018 323.266L41.5 315.086V55.7471L49.0693 48.6299L56.6514 41.5H231.391Z",
      glow: [
        { blur: 2, color: [0.482353, 0.219608, 0.101961], opacity: 0.4 },
        { blur: 8, color: [0.721569, 0.466667, 0.266667], opacity: 0.5 },
        { blur: 18, color: [0.957144, 0.852495, 0.564709], opacity: 0.2 },
      ],
      gradientStops: [
        { offset: 0, color: "#CE8819" },
        { offset: 0.983218, color: "#E17F43" },
      ],
      gradientCoords: { x1: 247.221, y1: 150.313, x2: 39.0114, y2: 278.836 },
    },
  },
};
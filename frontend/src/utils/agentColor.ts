const AGENT_COLORS = [
  "#FFB800",
  "#C87BFF",
  "#FF3B6B",
  "#39FF6D",
  "#FFCC33",
  "#6756FF",
  "#FF6B35",
  "#00E5A0",
  "#F43FFF",
  "#B6FF00", 
] as const;

export function getAgentColorByIndex(index: number) {
  return AGENT_COLORS[index % AGENT_COLORS.length];
}
const AGENT_COLORS = [
  "#FFB800",
  "#C87BFF",
  "#39FF6D",
  "#FF3B6B",
  "#6756FF",
  "#FF6B35",
  "#00E5A0",
] as const;

export function getAgentColorByIndex(index: number) {
  return AGENT_COLORS[index % AGENT_COLORS.length];
}
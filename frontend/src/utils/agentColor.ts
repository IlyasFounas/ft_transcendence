const AGENT_COLORS = [
  "#22D3EE",
  "#FFB800",
  "#C87BFF",
  "#FF3B6B",
  "#39FF6D",
  "#FFCC33",
  "#6FFFF6",
  "#6756FF",
] as const;

export function getAgentColorByIndex(index: number) {
  return AGENT_COLORS[index % AGENT_COLORS.length];
}
export const PODIUM_CARD_VARIANTS: Record<"gold" | "silver" | "bronze", {
  textColor: string;
  avatarSize: number;
  rankSize: string;
}> = {
  gold:   { textColor: "#DDB74E", avatarSize: 85, rankSize: "text-heading-large" },
  silver: { textColor: "#A8B4C0", avatarSize: 80, rankSize: "text-heading-medium" },
  bronze: { textColor: "#CE7B4A", avatarSize: 80, rankSize: "text-heading-medium" },
};
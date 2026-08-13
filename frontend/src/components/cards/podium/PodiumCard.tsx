// PodiumCard.tsx
import { PodiumShapeSvg } from "./PodiumShape";
import { AvatarFrame } from "./AvatarFrame";
import { PODIUM_VARIANTS } from "./PodiumShapeConfig";
import { PODIUM_CARD_VARIANTS } from "./PodiumCardConfig";

type PodiumCardProps = {
  variant: "gold" | "silver" | "bronze";
  rank: number;
  name: string;
  role: string;
  score: number;
  imageSrc: string;
  avatarSize?: number;
};

export default function PodiumCard({
  variant, rank, name, role, score, imageSrc,
}: PodiumCardProps) {
  const { width, height } = PODIUM_VARIANTS[variant];
  const { textColor, avatarSize, rankSize } = PODIUM_CARD_VARIANTS[variant];

  return (
    <div className="relative" style={{ width, height }}>
      <PodiumShapeSvg variant={variant} className="absolute inset-0" />

      <div
        className="absolute inset-0 flex justify-center items-center flex-col uppercase gap-2"
        style={{ color: textColor }}
      >
        <div className="text-center">
          <span className="text-hero">{rank}</span>
          <AvatarFrame variant={variant} size={avatarSize} imageSrc={imageSrc} />
        </div>
        <span className="text-heading-medium text-text-primary">{name}</span>
        <span className="text-body-primary">{role}</span>
        <div className="h-px w-[117px] bg-gradient-to-r from-transparent via-current to-transparent" />
        <div className="flex flex-row gap-[7px]">
          <span className={rankSize}>{score.toLocaleString()}</span>
          <span className="text-body-primary flex items-center">pts</span>
        </div>
      </div>
    </div>
  );
}
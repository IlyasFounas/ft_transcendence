import { BronzePodium, SilverPodium, GoldPodium } from "./PodiumShape";
import { BronzeAvatarHexagon, SilverAvatarHexagon, GoldAvatarHexagon } from "./HexagonShape";

interface PodiumProps {
	color: "bronze" | "silver" | "gold";
	pseudo: string;
	title: string;
	points: number;
}

const PODIUM_VARIANTS = {
  bronze: {
    Podium: BronzePodium,
    Avatar: BronzeAvatarHexagon,
    width: 287,
    height: 372,
	textColor: "#CE7B4A",
	rankSize: "text-heading-medium",
	rank: 3,
  },
  silver: {
    Podium: SilverPodium,
    Avatar: SilverAvatarHexagon,
    width: 287,
    height: 372,
	textColor: "#A8B4C0",
	rankSize: "text-heading-medium",
	rank: 2,
  },
  gold: {
    Podium: GoldPodium,
    Avatar: GoldAvatarHexagon,
    width: 312,
    height: 422,
	textColor: "#DDB74E",
	rankSize: "text-heading-large",
	rank: 1,
  },
};

export default function Podium({color, pseudo, title, points}: PodiumProps) {

  	const { Podium, Avatar, width, height, textColor, rankSize, rank } = PODIUM_VARIANTS[color];

	return (
		<div className="relative" style={{ width, height }}>
			<Podium className="absolute inset-0"></Podium>
			<div className="absolute inset-0 flex justify-center items-center flex-col uppercase gap-2"
				style={{color: textColor}}>
				<div className="text-center">
          			<span className="text-hero">{rank}</span>
					<Avatar></Avatar>
        		</div>
				<div className="flex flex-col gap-1 text-center">
        			<span className="text-heading-medium text-text-primary">{pseudo}</span>
        			<span className="text-body-primary">{title}</span>
        			<div className="h-px w-[117px] bg-gradient-to-r from-transparent via-current to-transparent" />
        			<div className="flex flex-row gap-[7px] justify-center">
          				<span className={rankSize}>{points}</span>
          				<span className="text-body-primary flex items-center">pts</span>
					</div>
        		</div>
			</div>
		</div>
	)
}
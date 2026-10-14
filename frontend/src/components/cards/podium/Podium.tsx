import { BronzePodium, SilverPodium, GoldPodium } from "./PodiumShape";
import {
	BronzeAvatarHexagon,
	SilverAvatarHexagon,
	GoldAvatarHexagon,
} from "./HexagonShape";

interface PodiumProps {
	color: "bronze" | "silver" | "gold";
	username: string;
	rank: string;
	avatar: string;
	points: number;
}

function fluidText(
	fontPx: number,
	lineHeightPx: number,
	fontWeight: number,
	letterSpacingEm: number,
	containerWidth: number,
	minRatio = 0.6
) {
	const idealCqw = (fontPx / containerWidth) * 100;
	const minPx = fontPx * minRatio;
	return {
		fontFamily: "var(--font-main)",
		fontWeight,
		letterSpacing: `${letterSpacingEm}em`,
		lineHeight: lineHeightPx / fontPx,
		fontSize: `clamp(${minPx}px, ${idealCqw}cqw, ${fontPx}px)`,
	};
}

const PODIUM_VARIANTS = {
	bronze: {
		Podium: BronzePodium,
		Avatar: BronzeAvatarHexagon,
		width: 287,
		height: 372,
		textColor: "#CE7B4A",
		position: 3,
	},
	silver: {
		Podium: SilverPodium,
		Avatar: SilverAvatarHexagon,
		width: 287,
		height: 372,
		textColor: "#A8B4C0",
		position: 2,
	},
	gold: {
		Podium: GoldPodium,
		Avatar: GoldAvatarHexagon,
		width: 312,
		height: 422,
		textColor: "#DDB74E",
		position: 1,
	},
};

export default function Podium({
	color,
	username,
	rank,
	avatar,
	points,
}: PodiumProps) {
	const { Podium, Avatar, width, height, textColor, position } =
		PODIUM_VARIANTS[color];

	const rankStyle = fluidText(40, 56, 700, 0.08, width, 0.5);
	const usernameStyle = fluidText(22, 30, 600, 0.03, width, 0.4);
	const bodyStyle = fluidText(14, 20, 600, 0.02, width, 0.5);

	const dividerWidthPx = 117;
	const idealDividerCqw = (dividerWidthPx / width) * 100;
	const dividerStyle = {
		width: `clamp(${dividerWidthPx * 0.6}px, ${idealDividerCqw}cqw, ${dividerWidthPx}px)`,
	};
	const pointsStyle =
		color === "gold"
			? fluidText(28, 35, 600, 0.05, width, 0.5)
			: usernameStyle;

	return (
		<div
			className="relative w-full @container"
			style={{
				aspectRatio: `${width} / ${height}`,
			}}
		>
			<Podium className="absolute inset-0 w-full h-full" />

			<div
				className="absolute inset-0 flex flex-col items-center justify-center gap-[3cqw] uppercase px-[4cqw]"
				style={{ color: textColor }}
			>
				<div className="w-full flex flex-col items-center justify-center gap-[2cqw]">
					<span style={rankStyle}>{position}</span>
					<Avatar avatar={avatar} className="w-[35%] h-auto" />
				</div>

				<div className="flex flex-col gap-[1.5cqw] text-center w-full">
					<span
						style={usernameStyle}
						className="text-text-primary truncate w-full px-6 md:px-10"
					>
						{username || "\u00A0"}
					</span>
					<span
						style={bodyStyle}
						className="text-text-primary truncate w-full px-6 md:px-10"
					>
						{rank || "\u00A0"}
					</span>

					<div
						className="h-px bg-gradient-to-r from-transparent via-current to-transparent mx-auto"
						style={dividerStyle}
					/>

					<div className="flex flex-row justify-center items-center gap-[2cqw]">
						<span style={pointsStyle}>{points}</span>
						<span style={bodyStyle} className="flex items-center">
							pts
						</span>
					</div>
				</div>
			</div>
		</div>
	);
}
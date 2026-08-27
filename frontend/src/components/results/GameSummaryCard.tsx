import Avatar from "../primitives/Avatar";

interface Player {
  username: string;
  role: string;
  avatar: string;
}
interface RoundResult {
  round: number;
  points: number;
}
interface GameSummaryCardProps {
  rank: number;
  rounds: RoundResult[];
  player: Player;
  displayRounds: number[];
}

const RANK_TEXT_COLORS = [
  "text-cyan-400",
  "text-amber-400",
  "text-purple-400",
  "text-rose-400",
  "text-emerald-400",
];

const RANK_HEX_COLORS = [
  "#22d3ee",
  "#fbbf24",
  "#c084fc",
  "#fb7185",
  "#34d399",
];

function getRankColors(rank: number) {
  const index = (rank - 1) % RANK_TEXT_COLORS.length;
  return {
    text: RANK_TEXT_COLORS[index],
    hex: RANK_HEX_COLORS[index],
  };
}

export const GRID_COLS =
  "grid-cols-[60px_minmax(0,1.5fr)_minmax(0,1fr)] " +
  "@[456px]:grid-cols-[60px_minmax(0,1.5fr)_minmax(0,1fr)_minmax(0,1fr)] " +
  "@[736px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(2,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[996px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(3,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1256px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(4,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1516px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(5,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1776px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(6,minmax(0,1fr))_minmax(0,1fr)]";

export const ROUND_VISIBILITY = [
  "hidden @[1776px]:flex",
  "hidden @[1516px]:flex",
  "hidden @[1256px]:flex",
  "hidden @[996px]:flex",
  "hidden @[736px]:flex",
  "hidden @[456px]:flex",
];

export default function GameSummaryCard({
  rank = 2,
  rounds,
  player,
  displayRounds,
}: GameSummaryCardProps) {
  	const offset = ROUND_VISIBILITY.length - displayRounds.length;
	const { text: rankColor, hex: rankHex } = getRankColors(rank);

  return (
    <div
      className={`w-full h-[65px] outline outline-1 outline-offset-[-1px] outline-cyan-800/30
        grid ${GRID_COLS} items-center px-4 uppercase ${rankColor}`}
    >
      <div className="text-center text-heading-medium">{rank}</div>

      <div className="flex items-center gap-[12px] @[640px]:px-[30px] min-w-0">
        <Avatar avatar={player.avatar} borderColor={rankHex} size={45} />
        <div className="flex flex-col min-w-0">
          <span className="text-heading-small truncate text-text-primary">
            {player.username}
          </span>
          <span className="hidden @[860px]:block text-body-primary truncate">
            {player.role}
          </span>
        </div>
      </div>

      {displayRounds.map((roundNumber, index) => {
        const round = rounds.find((r) => r.round === roundNumber);
        return (
          <div
            key={roundNumber}
            className={`${ROUND_VISIBILITY[offset + index]} items-center justify-center gap-[10px]`}
          >
            <span className="text-heading-small">
              {round ? round.points : "-"}
            </span>
            {round && <span className="text-body-primary">pts</span>}
          </div>
        );
      })}

      <div className="flex items-center justify-center gap-[10px]">
        <span className="text-heading-small">
          {rounds.reduce((total, r) => total + r.points, 0)}
        </span>
        <span className="text-body-primary">pts</span>
      </div>
    </div>
  );
}
import Avatar from "../primitives/Avatar";
import type { Player, RoundResult } from "../../types/game";
import { getAgentColorByIndex } from "../../utils/agentColor";

type Props = {
    rank: number;
    player: Player;
    rounds: RoundResult[];
    currentUsername: string;
    onMouseEnter?: () => void;
    onMouseLeave?: () => void;
};

export default function RankResultRow({
    rank = 2,
    rounds,
    player,
    currentUsername,
    onMouseEnter, 
    onMouseLeave
}: Props) {

    const isCurrentPlayer = player.username === currentUsername;
    const agentColor = isCurrentPlayer ? "#35DFFF" : getAgentColorByIndex(rank - 1);
    const totalPoints = rounds.reduce((total, r) => total + r.points, 0);

    return (
        <div
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
            className={`
                w-full 
                h-[65px] 
                outline outline-1 outline-offset-[-1px] 
                grid grid-cols-[90px_minmax(0,400px)_1fr]
                items-center 
                px-4 
                uppercase 
                hover:bg-cyan-400/8 transition-colors
                ${
                    isCurrentPlayer
                    ? "outline-cyan-400/60 bg-gradient-to-r from-cyan-400/15 via-cyan-400/5 to-transparent shadow-[0_0_20px_rgba(53,223,255,0.25)]"
                    : "outline-cyan-800/30"
                }`}
            style={{ color: agentColor }}
        >
            <div className="text-center flex  text-heading-medium">{rank}</div>

            <div className="flex items-center gap-[12px] min-w-0">
                <Avatar avatar={player.avatar} borderColor={agentColor} size={45} />
                <div className="flex flex-col min-w-0">
                    <span className="text-heading-small truncate text-white">
                        {player.username}
                    </span>
                    <span className="hidden @[456px]:block text-body-primary truncate">
                        {player.role}
                    </span>
                </div>
            </div>

            <div className="flex items-center justify-end gap-[10px]">
                <span className="text-heading-small">{totalPoints}</span>
                <span className="text-body-primary">pts</span>
            </div>
        </div>
    );
}
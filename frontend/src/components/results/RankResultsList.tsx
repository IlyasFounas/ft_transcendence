import RankResultRow from "./RankResultRow";
import { getSortedPlayersWithRank } from "../../utils/ranking.ts";
import type { PlayerResult } from "../../types/game.ts";
import { useState, useEffect } from "react"

function useIsLargeScreen() {
    const [isLarge, setIsLarge] = useState(
        typeof window !== "undefined" ? window.innerWidth >= 1024 : true
    );

    useEffect(() => {
        const mediaQuery = window.matchMedia("(min-width: 1024px)");
        const handleChange = () => setIsLarge(mediaQuery.matches);

        handleChange();
        mediaQuery.addEventListener("change", handleChange);
        return () => mediaQuery.removeEventListener("change", handleChange);
    }, []);

    return isLarge;
}

interface RankResultListProps {
    players: PlayerResult[];
    currentUsername: string;
    columns?: 1 | 2;
    scrollable?: boolean;
    onHoverPlayer?: (username: string | null) => void;
}

export default function RankResultList({
    players,
    currentUsername,
    columns = 1,
    scrollable = false,
    onHoverPlayer
}: RankResultListProps) {

    const isLargeScreen = useIsLargeScreen();
    const sortedPlayers = getSortedPlayersWithRank(players);

    const shouldSplit = columns === 2 && isLargeScreen;

    const columnGroups = shouldSplit
        ? [
            sortedPlayers.slice(0, 5),
            sortedPlayers.slice(5, 10),
        ]
        : [sortedPlayers];

    return (
        <div className={`grid gap-4 w-full ${shouldSplit ? "lg:grid-cols-2" : "grid-cols-1"} flex flex-col flex-1 min-h-0 w-full`}>
            {columnGroups.map((group, groupIndex) => (
                <div
                    key={groupIndex}
                    className="@container flex flex-col rounded-[12px] bg-surface-base border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden"
                >
                    <div className="sticky top-0 z-10 h-[40px] grid grid-cols-[90px_minmax(0,400px)_1fr] items-center px-4
                        uppercase text-heading-small text-text-primary/70 bg-surface-base bg-linear-to-t from-black to-surface-base shrink-0">
                        <h1>Rank</h1>
                        <h1 className="px-8">Agent</h1>
                        <h1 className="flex justify-end">Total</h1>
                    </div>

                    <div className={scrollable ? "h-[230px] lg:h-[600px] overflow-y-auto [scrollbar-gutter:stable] flex-1 min-h-0 flex flex-col" : ""}>
                        {group.map((result) => (
                            <RankResultRow
                                key={result.player.username}
                                rank={result.rank}
                                player={result.player}
                                rounds={result.rounds}
                                currentUsername={currentUsername}
                                onMouseEnter={() => onHoverPlayer?.(result.player.username)}
                                onMouseLeave={() => onHoverPlayer?.(null)}
                            />
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
}
import type { RoundResult } from "../../types/game";
import GameResultRow from "./GameResultRow";

interface GameResultListProps {
  rounds: RoundResult[];
  currentRound: number;
  totalRounds: number;
}

export default function GameResultList({ currentRound, rounds, totalRounds }: GameResultListProps) {
  const displayRounds = Array.from({ length: totalRounds }, (_, i) => {
	const roundNb = i + 1;
	const played = rounds.find((r) => r.round === roundNb);
	return {
	  round: roundNb,
	  points: played?.points ?? null,
	  distance: played?.distance ?? null,
	};
  });

	return (
		<div className="flex flex-col flex-1 min-h-0 w-full rounded-[12px] border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden">
		
		<div className="flex flex-row shrink-0 items-center h-[40px] justify-around uppercase text-heading-small text-text-primary/70 bg-gradient-to-t from-black to-bg-surface-base">
			<h1>Round</h1>
			<h1>Score</h1>
			<h1>Distance</h1>
		</div>
		
		<div className="flex-1 flex flex-col min-h-0 overflow-y-auto [scrollbar-gutter:stable] bg-surface-base">
			{displayRounds.map((round) => (
			<GameResultRow
				currentRound={currentRound}
				key={round.round}
				round={round.round}
				distance={round.distance}
				points={round.points}
			/>
			))}
		</div>
		</div>
	);
}
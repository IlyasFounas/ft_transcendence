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
	<div className="flex flex-col rounded-[12px] h-full border border-border-default/30 shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden w-full max-h-[270px] lg:max-h-[542px]">
	  <div className="flex flex-row items-center h-[40px] justify-around uppercase text-heading-small text-text-primary/70 bg-surface-base ">
		<h1>Round</h1>
		<h1>Score</h1>
		<h1>Distance</h1>
	  </div>
	  <div className="max-h-[230px] overflow-y-auto lg:max-h-[500px]">
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
import { useState } from "react"
import SecondaryButton from "../components/buttons/SecondaryButton";
import { Users, Globe } from "lucide-react"
import Podium from "../components/cards/podium/Podium";
import { mockGameResult } from "../mocks/gameResult.mock";
import LeaderBoardList from "../components/results/LeaderboardList";
import { getSortedPlayersWithRank } from "../utils/ranking.ts";
import { getTotalPoints } from "../utils/getTotalPoints.ts";

export default function LeaderboardPage() {
	
	const [mode, setMode] = useState<"global" | "friends">("friends");
  	const gameResult = mockGameResult;
	const sortedPlayers = getSortedPlayersWithRank(gameResult.players);
	const topthree = sortedPlayers.slice(0, 3);
	const rest = sortedPlayers.slice(3);

	const currentPlayerResult = sortedPlayers.find(
  (result) => result.player.username === gameResult.currentUsername
);

	return (
		<div>
			<div className="
				absolute 
				inset-0 
				-z-10
				opacity-50
				bg-[url(/backgrounds/leaderboard_background.png)]
				h-screen">
			</div>
			<div className="
			border py-[24px] px-[50px]
			flex flex-row justify-between m-8
			bg-surface-base
			border-border-default rounded-xl
			">
				<div>
					<h1 className="
						text-hero
						uppercase
					"
					>{mode === "global" ? "Global Agent Ranking" : "Friends Agent Ranking"}</h1>
					<p className="
						text-label-small
						text-text-secondary
					">
						{mode === "global" 
						? "Top agents across all the operation worldwide." 
						: "Top agents among your friends."}
					</p>
				</div>
			<div className="
			flex gap-[30px]">
				<SecondaryButton
					active={mode === "global"}
					icon={<Globe/>}
					onClick={() => setMode("global")}
					>
					Global
				</SecondaryButton>
				<SecondaryButton
					active={mode === "friends"}
					icon={<Users/>}
					onClick={() => setMode("friends")}
					>
					Friends
				</SecondaryButton>
			</div>
		</div>
		<div className="flex flex-row justify-center gap-20 items-start">
			<div className="grid grid-cols-[1fr_1.1fr_1fr] w-1/2 px-12">
				<div className="mt-40 -mx-8">
					<Podium color="silver" username={topthree[1].player.username} title={topthree[1].player.role} points={getTotalPoints(topthree[1].rounds)}  />
				</div>
				<div className="-mx-8">
					<Podium color="gold" username={topthree[0].player.username} title={topthree[0].player.role} points={getTotalPoints(topthree[0].rounds)} />
				</div>
				<div className="mt-40 -mx-8">
					<Podium color="bronze" username={topthree[2].player.username} title={topthree[2].player.role} points={getTotalPoints(topthree[2].rounds)} />
				</div>
			</div>
			<div className="flex justify-center">
				<LeaderBoardList
					players={rest}
					currentPlayerResult={currentPlayerResult}
					currentUsername={gameResult.currentUsername}
				/>
			</div>
		</div>
	</div>
 	);
}
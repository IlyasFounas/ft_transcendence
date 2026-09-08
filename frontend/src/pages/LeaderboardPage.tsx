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
			border py-[12px] md:py-[24px] px-[10px] 
			flex md:flex-row 
			flex-col justify-center text-center md:text-left gap-4 md:justify-between mx-8 mt-8 md:mb-8
			bg-surface-base
			border-border-default rounded-xl items-center
			">
				<div>
					<h1 className="
						text-heading-medium lg:text-hero
						uppercase
					"
					>{mode === "global" ? "Global Agent Ranking" : "Friends Agent Ranking"}</h1>
					<p className="
						text-label-small
						text-text-secondary
						uppercase
					">
						{mode === "global" 
						? "Top agents across all the operation worldwide." 
						: "Top agents among your friends."}
					</p>
				</div>
				<div className="flex justify-center gap-4 md:gap-[30px] h-[60px]">
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
		<div className="flex lg:flex-row flex-col justify-center ">
			<div className="grid grid-cols-2 sm:grid-cols-[1fr_1.1fr_1fr] gap-0 lg:w-1/2 items-start">
				<div className="col-span-2 w-[120px] sm:w-full mx-auto flex justify-center sm:col-span-1 sm:order-2">
					<Podium
					color="gold"
					username={topthree[0].player.username}
					title={topthree[0].player.role}
					points={getTotalPoints(topthree[0].rounds)}
					/>
				</div>
				<div className="ml-auto sm:order-1 w-[120px] sm:w-full sm:mt-40">
					<Podium
					color="silver"
					username={topthree[1].player.username}
					title={topthree[1].player.role}
					points={getTotalPoints(topthree[1].rounds)}
					/>
				</div>
				<div className="sm:order-3 mr-auto w-[120px] sm:w-full sm:mt-40">
					<Podium
					color="bronze"
					username={topthree[2].player.username}
					title={topthree[2].player.role}
					points={getTotalPoints(topthree[2].rounds)}
					/>
				</div>
				</div>
			<div className="px-[20px] md:px-[50px] pb-10 lg:w-1/2">
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
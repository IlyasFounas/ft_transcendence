import { useState, useEffect } from "react"
import SecondaryButton from "../components/buttons/SecondaryButton";
import { Users, Globe } from "lucide-react"
import Podium from "../components/cards/podium/Podium";
import LeaderBoardList from "../components/leaderboard/LeaderboardList";
import { apiFetch } from "../api/apiFetch";
import { getGlobalLeaderboard, type LeaderboardEntry } from "../services/leaderboard";
import { getFriends, type FriendsEntry } from "../services/friends";
import { useAuth } from "../context/AuthContext";

export default function LeaderboardPage() {
	const [mode, setMode] = useState<"global" | "friends">("friends");
	const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
	const [friends, setFriends] = useState<FriendsEntry[]>([]);
	const { user, authFetch } = useAuth();

	useEffect(() => {
		getGlobalLeaderboard(apiFetch).then(setLeaderboard);
		getFriends(authFetch).then(setFriends);
	}, []);

	const friendUsernames = new Set(friends.map((f) => f.name));
	friendUsernames.add(user?.Username ?? "");

	const displayedLeaderboard =
		mode === "friends"
		? leaderboard
			.filter((entry) => friendUsernames.has(entry.username))
			.map((entry, i) => ({ ...entry, position: i + 1 }))
		: leaderboard;

	const topthree = displayedLeaderboard.slice(0, 3);
	const rest = displayedLeaderboard.slice(3);

	const currentPlayerResult = displayedLeaderboard.find(
		(entry) => entry.username === user?.Username
	);

	return (
		<div>
		<div className="
					absolute 
					inset-0 
					-z-10
					opacity-50
					bg-[url(/backgrounds/leaderboard_background.png)]
					bg-cover bg-center bg-no-repeat
					h-screen">
		</div>
		<div className="
				border py-[12px] md:py-[24px] px-8
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
				username={topthree[0]?.username ?? ""}
				rank={topthree[0]?.rank ?? ""}
				avatar={topthree[0]?.avatar ?? ""}
				points={topthree[0]?.total ?? 0}
				/>
			</div>
			<div className="ml-auto sm:order-1 w-[120px] sm:w-full sm:mt-40">
				<Podium
				color="silver"
				username={topthree[1]?.username ?? ""}
				rank={topthree[1]?.rank ?? ""}
				avatar={topthree[1]?.avatar ?? ""}
				points={topthree[1]?.total ?? 0}
				/>
			</div>
			<div className="sm:order-3 mr-auto w-[120px] sm:w-full sm:mt-40">
				<Podium
				color="bronze"
				username={topthree[2]?.username ?? ""}
				rank={topthree[2]?.rank ?? ""}
				avatar={topthree[2]?.avatar ?? ""}
				points={topthree[2]?.total ?? 0}
				/>
			</div>
			</div>
			<div className="px-[20px] md:px-[50px] pb-10 lg:w-1/2">
			<LeaderBoardList
				players={rest}
				currentPlayerResult={currentPlayerResult}
				currentUsername={user?.Username ?? ""}
				totalRank={mode === "global" ? 100 : 10}
			/>
			</div>
		</div>
		</div>
	);
}
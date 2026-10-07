import { useAuth } from "../context/AuthContext";
import Avatar from "../components/primitives/Avatar";
import { Pencil, MapPin, ClipboardList, Crosshair, Trophy } from "lucide-react";
import CalculateRank from "../utils/rank";
import SecondaryButton from "../components/buttons/SecondaryButton";
import StatCard from "../components/cards/StatCard";
import OperationHistory from "../components/cards/OperationHistory";
import FriendListPanel from "../components/cards/FriendListPanel";
import { useState, useEffect } from "react";
import EditProfileModal from "../components/cards/EditProfileModal";
import type { UserStatistics } from "../types/statistics.type";
import type { UserLeaderboard } from "../types/user.types";
import type { Operation } from "../types/operation.types";
import { formatOperationDate } from "../utils/formatDate";

interface HistoricGameResponse {
    id: number;
    xp: number;
    date: string;
}

export default function ProfilePage() {

	const { user, authFetch } = useAuth();
    const [statistics, setStatistics] = useState<UserStatistics | null>(null);
    const [leaderboard, setLeaderboard] = useState<UserLeaderboard[]>([]);
    const [operations, setOperations] = useState<Operation[]>([]);

	const totalScore = user?.XP ?? 0;
	const rank = CalculateRank(totalScore);

    useEffect(() => {

        const getStatistics = async () => {
            const result = await authFetch("api/user/statistics");
            if (!result) {
                return ;
            }
            
            setStatistics(result.data.statistics);
        };

        const getLeaderboard = async () => {
            const result = await authFetch("api/leaderboard/user");

            if (!result) {
                return ;
            }
            setLeaderboard(result.data.leaderboard);
        };

        const getOperations = async () => {
        if (!user?.ID) {
            return;
        }

        const result = await authFetch(`api/game/historic/${user.ID}`);

        if (!result || result.status !== 200) {
            return;
        }

        const historicGames: HistoricGameResponse[] = result.data.historic_games;

        setOperations(
            historicGames.map((game) => ({
                operationId: game.id.toString(),
                points: game.xp,
                dateTime: formatOperationDate(game.date),
            }))
        );
    };

        getStatistics();
        getLeaderboard();
        getOperations();
    }, []);


    const [showEditProfile, setShowEditProfile] = useState(false);

    return (
    	<main
			className="
			flex-1
			min-h-0
			w-full
			bg-surface-base
			px-8
			py-5
			lg:px-12
			overflow-y-auto
            xl:gap-5
            xl:flex
            xl:flex-col
          	"
      	>
			<section
                className="
                    w-full
                    border
                    border-border-default
                    rounded-md
                    bg-surface-base
                    p-5
                    md:p-6
                    xl:p-10
                    flex
                    flex-col
                    xl:flex-row
                    gap-4
                    lg:gap-10
                    xl:gap-15
                    items-center
                    xl:items-start
                    xl:shrink-0
                "
            >
				<div className="shrink-0 md:hidden">
                    <Avatar
                        avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                        size={80}
                    />
                </div>
                <div className="shrink-0 hidden md:block xl:hidden">
                    <Avatar
                        avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                        size={120}
                    />
                </div>
                <div className="shrink-0 hidden xl:block">
                    <Avatar
                        avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                        size={200}
                    />
                </div>

				<div className="flex-1 w-full flex flex-col gap-3 xl:gap-8 text-center xl:text-left xl:mt-5">

					<div>
                        <h1 className="text-heading-medium md:text-heading-large xl:text-hero font-semibold text-text-primary uppercase">
                            Agent {user?.Username ?? "Player"}
                        </h1>
                        <p className="text-label-small md:text-heading-small lg:text-heading-medium text-primary uppercase">
                            {user?.Rank ?? "Junior Agent"}
                        </p>
                    </div>
					<div className="w-full xl:max-w-xl">

                        <div className="flex justify-between text-label-small md:text-body-primary lg:text-heading-small xl:text-heading-medium text-text-secondary uppercase mb-1">
							<span>Next Rank</span>
							<span>
								{totalScore.toLocaleString()} / {(totalScore + rank.xpRemaining).toLocaleString()} XP
							</span>
						</div>

						<p className="text-label-small md:text-body-primary lg:text-heading-small xl:text-heading-medium text-left text-text-primary uppercase mb-2">
							{rank.nextRank}
						</p>

                        <div className="h-2 w-full bg-surface-elevated rounded-full overflow-hidden">
                            <div
                                className="h-full bg-primary rounded-full"
                                style={{ width: `${rank.progress}%` }}
                            />
                        </div>
                    </div>
				</div>
				<div className="flex items-center justify-center text-center xl:self-center xl:gap-20">

					<div className="text-center uppercase hidden md:block">
						<p className="text-label-small md:text-heading-small xl:text-heading-medium text-text-secondary">
							Global Rank
						</p>

						<p className="text-heading-medium xl:text-hero text-primary font-bold">
							{leaderboard[0]?.position ?? ""}
						</p>
					</div>
					<div className="xl:pt-5">
						<SecondaryButton 
                            className="text-heading-small md:text-heading-medium lg:text-heading-large scale-60! xs:hover:scale-60! md:scale-80 xl:scale-100 md:hover:scale-80 xl:hover:scale-100"
                            onClick={() => setShowEditProfile(true)}
                        >
							<Pencil size={18} />
							Edit Profile
						</SecondaryButton>
					</div>

				</div>

			</section>

			<section
                className="
                    mt-2
                    border
                    border-border-default
                    rounded-md
                    p-4
                    xl:p-10
                    xl:shrink-0
                "
            >

                <h2 className="text-center text-heading-small xl:text-heading-large text-text-primary uppercase mb-4">
                    Statistics
                </h2>

                <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-6">

                    <StatCard
                        title="Operations"
                        value={statistics?.finished_games?.toString() ?? "0"}
                        subtitle="Total Completed"
                        icon={<ClipboardList size={60} strokeWidth={1.5} />}
                    />
                    <StatCard
                        title="Accuracy"
                        value={`${statistics?.accuracy ?? 0}%`}
                        subtitle="Across All Operations"
                        icon={<Crosshair size={60} strokeWidth={1.5} className="text-status-green"/>}
                    />
                    <StatCard
                        title="Best Score"
                        value={statistics?.best_score?.toLocaleString() ?? "0"}
                        subtitle="Single Best Game"
                        icon={<Trophy size={60} strokeWidth={1.5} className="text-status-yellow"/>}
                    />
                    <StatCard
                        title="Closest Guess"
                        value={
                            statistics?.closest_guess
                                ? `${statistics.closest_guess.value} ${statistics.closest_guess.unit.toUpperCase()}`
                                : "N/A"
                        }
                        subtitle="Best Distance To Target"
					icon={<MapPin size={60} strokeWidth={1.5} className="text-[#C87BFF]"/>}
                    />

                </div>

            </section>

			<section
                className="
                    mt-2
                    grid
                    grid-cols-1
                    xl:grid-cols-2
                    gap-2
                    xl:gap-8
                    xl:h-[calc(100vh-700px)]
                    xl:min-h-100
                "
            >
				<OperationHistory operations={operations} />
				<FriendListPanel />
			</section>

            {showEditProfile && (
                <EditProfileModal
                    username={user?.Username ?? ""}
                    avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                    avatars={[
                        "/avatars/avatar1.png",
                        "/avatars/avatar2.png",
                        "/avatars/avatar3.png",
                        "/avatars/avatar4.png",
                        "/avatars/avatar5.png",
                        "/avatars/avatar6.png",
                        "/avatars/avatar7.png",
                        "/avatars/avatar8.png",
                        "/avatars/avatar9.png",
                        "/avatars/avatar10.png",
                        "/avatars/avatar11.png",
                        "/avatars/avatar12.png",
                        "/avatars/avatar13.png",
                        "/avatars/avatar14.png",
                        "/avatars/avatar15.png",
                        "/avatars/avatar16.png",
                        "/avatars/avatar17.png",
                        "/avatars/avatar18.png",
                    ]}
                    onClose={() => setShowEditProfile(false)}
                />
            )}
      	</main>					    
  	);
}
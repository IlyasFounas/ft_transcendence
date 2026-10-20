import PrimaryButton from "../components/buttons/PrimaryButton";
import ProfileCard from "../components/cards/ProfileCard";
import { useAuth } from "../context/AuthContext";
import CalculateRank from "../utils/rank";
import { useEffect, useState } from "react";
import type { UserStatistics } from "../types/statistics.type";
import { useNavigate } from "react-router-dom";
import { useLobbySocket } from "../context/GameSocketContext";

export default function HomePage() {
    const { user, authFetch, getUser } = useAuth();
    const [statistics, setStatistics] = useState<UserStatistics | null>(null);
	const { connect } = useLobbySocket();

    const totalScore = user?.XP ?? 0;
    const rank = CalculateRank(totalScore);
    
    useEffect(() => {
        const refreshUser = async () => {
            await getUser();
        };

        const getStatistics = async () => {
            const result = await authFetch("api/user/statistics");

            if (!result) {
                return ;
            }
            
            setStatistics(result.data.statistics);
        };
        
        refreshUser();
        getStatistics();
    }, []);

	const handlePlaySolo = async () => {

		const current = await authFetch("/api/lobby/current");
		if (current?.status === 200 && current.data?.game_id) {
			const modeRes = await authFetch(`/api/game/${current.data.game_id}/mode`);
			connect();
			if (modeRes?.data?.mode === "multi") {
				navigate("/lobby");
			} else {
				navigate(`/game/${current.data.game_id}`);
			}
			return;
		}
		
		const body = new URLSearchParams({
				timer: "120",
				round_total: "5", 
			});

		const result = await authFetch("/api/lobby/create", {
			method: "POST",
			headers: { "Content-Type": "application/x-www-form-urlencoded" },
			body: body.toString(),
		});

		if (result?.status === 200) {
			const cur = await authFetch("/api/lobby/current");
			const gameId = cur?.data?.game_id;
			if (gameId) {
				await authFetch(`/api/game/${gameId}/mode`, {
					method: "PATCH",
					headers: { "Content-Type": "application/x-www-form-urlencoded" },
					body: new URLSearchParams({ mode: "solo" }).toString(),
				});
			}
			connect();
			navigate("/game/solo");
		}
	};

	const navigate = useNavigate();

    return (
        
        <main
            className="
                flex-1
                w-full
                min-w-0
                bg-surface-base
                bg-cover
                bg-center
                bg-[url(/backgrounds/home_background.png)]
                bg-no-repeat
                flex
                flex-col
                items-center
                justify-start
                xl:items-start
                px-5
                py-6
                md:py-8
                xl:px-20
                gap-6
                md:gap-10
                "
        >
            <div>
                <ProfileCard 
                    username={user?.Username ?? "player"}
                    avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                    role={user?.Rank ?? "Junior Agent"}
                    operations={statistics?.finished_games ?? 0}
                    accuracy={statistics?.accuracy ?? 0}
                    score={statistics?.best_score ?? 0}
                    nextRank={rank.nextRank ?? ""}
                    xpToNextRank={rank.xpRemaining}
                    progress={rank.progress}
                />
            </div>
            <div className="
                flex
                flex-col
                mt-auto
                
                items-center
                xl:self-center

                space-y-2
                lg:space-y-6
                lg:px-20
                lg:pb-12
                xl:pr-0
                "
            >
                <div className="
                    flex
                    w-62.5
                    gap-4
                    items-center
                    sm:px-6
                    md:w-105
                    lg:w-125
                    "
                >
                    <hr className="h-px flex-1 border-primary" />
                    <span className="text-body-primary md:text-heading-small lg:text-heading-medium text-primary uppercase">
                        your mission awaits
                    </span>
                    <hr className="h-px flex-1 border-primary" />
                </div>
                <div className="flex gap-3 lg:gap-8">
                    <p className="text-heading-small md:text-heading-medium lg:text-hero font-semibold text-text-primary">ANALYZE.</p>
                    <p className="text-heading-small md:text-heading-medium lg:text-hero font-semibold text-text-primary">LOCATE.</p>
                    <p className="text-heading-small md:text-heading-medium lg:text-hero font-semibold text-text-primary">REPORT.</p>
                </div>
                <div className="w-full md:hidden">
                    <PrimaryButton size="small" subtitle="Launch a new operation" onClick={handlePlaySolo}>
                        Start Mission
                    </PrimaryButton>
                </div>
                <div className="w-full hidden md:block lg:hidden text-center">
                    <PrimaryButton size="medium" subtitle="Launch a new operation" onClick={handlePlaySolo}>
                        Start Mission
                    </PrimaryButton>
                </div>    
                <div className="w-full hidden lg:block">
                    <PrimaryButton size="large" subtitle="Launch a new operation" onClick={handlePlaySolo}>
                        Start Mission
                    </PrimaryButton>
                </div>              
            </div>
        </main>
    );
}

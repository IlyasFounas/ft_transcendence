import { useAuth } from "../context/AuthContext";
import Avatar from "../components/primitives/Avatar";
import { Pencil, MapPin, ClipboardList, Crosshair, Trophy } from "lucide-react";
import CalculateRank from "../utils/rank";
import SecondaryButton from "../components/buttons/SecondaryButton";
import StatCard from "../components/cards/StatCard";
import OperationHistory from "../components/cards/OperationHistory";
import { operations } from "../data/operations";
import FriendListPanel from "../components/cards/FriendListPanel";
import { friends, requests, pending } from "../data/friends";

export default function ProfilePage() {

	const { user } = useAuth();

	const totalScore = 24658;
	const rank = CalculateRank(totalScore);

    return (
    	<main
			className="
			flex-1
			min-h-0
			w-full
			bg-surface-base
			px-5
			py-5
			lg:px-12
			overflow-y-auto
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
                    lg:p-6
                    flex
                    flex-col
                    lg:flex-row
                    gap-6
                    lg:gap-10
                    items-center
                    lg:items-start
                "
            >
				<div className="shrink-0">
                    <Avatar
                        avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                        size={200}
                    />
                </div>

				<div className="flex-1 w-full flex flex-col gap-3 lg:mt-5">

					<div>
                        <h1 className="text-hero font-semibold text-text-primary uppercase">
                            Agent {user?.Username ?? "Player"}
                        </h1>
                        <p className="text-heading-medium text-primary uppercase">
                            {rank.currentRank}
                        </p>
                    </div>
					<div className="w-full max-w-xl">

                        <div className="flex justify-between text-heading-small text-text-secondary uppercase mb-1">
							<span>Next Rank</span>
							<span>
								{totalScore.toLocaleString()} / {(totalScore + rank.xpRemaining).toLocaleString()} XP
							</span>
						</div>

						<p className="text-heading-small text-text-primary uppercase mb-2">
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
				<div className="flex items-center justify-center text-center gap-20 lg:self-center lg:mr-10">

					<div className="text-center uppercase">
						<p className="text-heading-medium text-text-secondary">
							Global Rank
						</p>

						<p className="text-[48px] text-primary font-bold">
							#18
						</p>
					</div>
					<div className="pt-5">
						<SecondaryButton>
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
                    lg:p-5
                "
            >

                <h2 className="text-center text-heading-medium text-text-primary uppercase mb-4">
                    Statistics
                </h2>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">

                    <StatCard
                        title="Operations"
                        value="152"
                        subtitle="Total Completed"
                        icon={<ClipboardList size={75} strokeWidth={1.5}  />}
                    />
                    <StatCard
                        title="Accuracy"
                        value="74%"
                        subtitle="Across All Operations"
                        icon={<Crosshair size={75} strokeWidth={1.5} className="text-status-green"/>}
                    />
                    <StatCard
                        title="Best Score"
                        value="24,580"
                        subtitle="Single Round Best"
                        icon={<Trophy size={75} strokeWidth={1.5} className="text-status-yellow"/>}
                    />
                    <StatCard
                        title="Closest Guess"
                        value="1.2 KM"
                        subtitle="Best Distance To Target"
					icon={<MapPin size={75} strokeWidth={1.5} className="text-[#C87BFF]"/>}
                    />

                </div>

            </section>

			<section
				className="
					mt-2
					grid
					grid-cols-1
					lg:grid-cols-2
					gap-8
				"
            >
				<OperationHistory operations={operations} />
				<FriendListPanel
					friends={friends}
					requests={requests}
					pending={pending}
				/>
			</section>
      	</main>					    
  	);
}		
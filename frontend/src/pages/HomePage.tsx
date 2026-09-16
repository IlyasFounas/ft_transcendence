import PrimaryButton from "../components/buttons/PrimaryButton";
import ProfileCard from "../components/cards/ProfileCard";
import { useAuth } from "../context/AuthContext";

export default function HomePage() {
    const { user } = useAuth();
    return (
        <main
            className="
                flex-1
                w-full
                bg-surface-base
                bg-cover
                bg-center
                bg-[url('/images/global-bg.png')]
                bg-no-repeat
                flex
                flex-col
                items-center
                justify-start
                px-5
                py-12
                lg:px-20
                lg:py-20
                lg:items-start
                gap-6
                md:gap-20
            "
        >
            <div>
                <ProfileCard 
                    username={user?.Username ?? "player"}
                    avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                    role="field analyst"
                    operations={123}
                    accuracy={74}
                    score={24658}
                    nextRank="Senior Analyst"
                    xpToNextRank={2940}
                />
            </div>
            <div className="flex flex-col space-y-2 lg:space-y-6 items-center lg:self-center mt-auto lg:pb-12">
                <div className="flex w-[250px] lg:w-[500px] gap-4 items-center sm:px-6">
                    <hr className="h-px flex-1 border-primary" />
                    <span className="text-body-primary lg:text-heading-medium text-primary uppercase">
                        your mission awaits
                    </span>
                    <hr className="h-px flex-1 border-primary" />
                </div>
                <div className="flex gap-3 lg:gap-8">
                    <p className="text-heading-small lg:text-hero font-semibold text-text-primary">ANALYZE.</p>
                    <p className="text-heading-small lg:text-hero font-semibold text-text-primary">LOCATE.</p>
                    <p className="text-heading-small lg:text-hero font-semibold text-text-primary">REPORT.</p>
                </div>
                <PrimaryButton size="small" subtitle="Launch a new operation">
                    Start Mission
                </PrimaryButton>
            </div>
        </main>
    );
}

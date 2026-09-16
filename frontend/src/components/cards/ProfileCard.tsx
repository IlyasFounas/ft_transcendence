import Avatar from "../primitives/Avatar";

interface ProfileCardProps {
    username: string;
    avatar: string;
    role: string;
    operations: number;
    accuracy: number;
    score: number;
    nextRank: string;
    xpToNextRank: number;
}

export default function ProfileCard({
    username,
    avatar,
    role,
    operations,
    accuracy,
    score,
    nextRank,
    xpToNextRank,
}: ProfileCardProps) {
    return (
        <div className="
                bg-surface-base
                border
                border-border-default
                rounded-2xl
                space-y-4
                p-6
                md:p-5
                w-[300px]
                lg:w-[480px]
            "
        >
            <div className="w-full flex flex-col lg:flex-1 lg:item lg:text-left lg:flex-row gap-2 lg:gap-6 items-center">
                <div className="shrink-0">
                    <Avatar avatar={avatar} size={100} />
                </div>
                <div className="
                    flex
                    flex-col
                    gap-1
                    w-full
                    items-center
                    lg:items-start
                    lg:min-w-0
                    lg:flex-1
                    "
                >
                    <div className="
                        truncate
                        min-w-0
                        w-full
                        text-center
                        lg:text-left
                        text-heading-small
                        lg:text-heading-large
                        text-text-primary
                        uppercase
                        "
                    >
                        <span>Agent </span>{username}
                    </div>
                    <div className="text-body-primary lg:text-heading-small text-primary uppercase">
                        {role}
                    </div>
                </div>
            </div>

            <div className="flex flex-row uppercase">
                <div className="flex-1">
                    <p className="text-label-small lg:text-body-primary text-text-secondary">Operation</p>
                    <p className="text-heading-small lg:text-heading-medium text-text-primary">{operations}</p>
                </div>
                <div className="flex-1">
                    <p className="text-label-small lg:text-body-primary text-text-secondary">Accuracy</p>
                    <span className="text-heading-small lg:text-heading-medium text-text-primary">{accuracy}</span><span> %</span>
                </div>
                <div className="flex-1">
                    <p className="text-label-small lg:text-body-primary text-text-secondary">Highest Score</p>
                    <p className="text-heading-small lg:text-heading-medium text-text-primary">{score}</p>
                </div>
            </div>

            <div className="flex flex-col uppercase space-y-1 lg:space-y-2">
                <p className="text-label-small lg:text-body-primary text-text-secondary">Next Rank</p>
                <p className="text-body-primary text-text-primary">{nextRank}</p>
                <div className="space-y-1">
                    <div className="h-1 bg-white rounded-2xl">
                        <div className="h-1 bg-primary rounded-2xl w-1/2"></div>
                    </div>
                    <span className="text-label-small text-text-secondary uppercase">
                        {xpToNextRank} XP to Next Rank
                    </span>
                </div>
            </div>

        </div>
    );
}
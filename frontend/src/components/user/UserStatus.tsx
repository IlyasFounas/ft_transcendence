import type { UserStatusType } from "../../types/user.types";

interface UserStatusProps {
    status: UserStatusType;
}

export default function UserStatus({
    status,
}: UserStatusProps) {
    const statusClasses = {
        online: "text-status-green",
        "in-mission": "text-status-yellow",
        offline: "text-text-secondary",
    };

    const statusLabels = {
        online: "ONLINE",
        "in-mission": "IN MISSION",
        offline: "OFFLINE",
    };

    return (
        <div
            className={`
                flex
                items-center
                gap-2
                text-heading-small
                ${statusClasses[status]}
            `}
        >
            <span className="w-2 h-2 rounded-full bg-current" />
            <span>{statusLabels[status]}</span>
        </div>
    );
}
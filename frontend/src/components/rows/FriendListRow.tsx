import UserProfile from "../user/UserProfile";
import SmallButton from "../buttons/SmallButton";
import UserStatus from "../user/UserStatus";
import type { UserProps, UserStatusType } from "../../types/user.types";

interface FriendListRowProps extends UserProps {
    status: UserStatusType;
    rank: number;
    onRemove?: () => void;
}

export default function FriendListRow({
    username,
    role,
    avatar,
    status,
    rank,
    onRemove,
}: FriendListRowProps) {
    return (
        <div className="
            grid
            grid-cols-[1.5fr_1fr_1fr_auto]
            items-center
            bg-surface-base
            border
            border-border-default
            px-4
            h-16
            "
        >
            <UserProfile 
                username={username}
                role={role}
                avatar={avatar}
            />
            <UserStatus status={status}/>
            <div className="text-center">
                <div className="text-heading-small text-text-primary">
                    #{rank}
                </div>
                <div className="text-caption-medium text-text-secondary uppercase">
                    Global Rank
                </div>
            </div>
            <SmallButton variant="danger" onClick={onRemove}>
                Remove
            </SmallButton>
        </div>
    );
}
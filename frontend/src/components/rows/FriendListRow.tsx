import UserProfile from "../user/UserProfile";
import SmallButton from "../buttons/SmallButton";
import type { UserProps } from "../../types/user.types";

interface FriendListRowProps extends UserProps {
    onRemove?: () => void;
}

export default function FriendListRow({
    username,
    role,
    avatar,
    onRemove,
}: FriendListRowProps) {
    return (
        // <div className="
        //     grid
        //     grid-cols-[1.5fr_1fr_1fr_auto]
        //     items-center
        //     bg-surface-base
        //     border
        //     border-border-default
        //     px-4
        //     h-16
        //     "
        // >
        <div className="
            flex
            items-center
            justify-between
            bg-surface-base
            border
            border-border-default
            px-4
            h-16
            xl:h-20
            xl:px-8
            "
        >
            <UserProfile 
                username={username}
                role={role}
                avatar={avatar}
            />
            {/* optional before link with websocket */}
            {/* {status && <UserStatus status={status} />}
            {rank !== undefined && (
                <div className="text-center">
                    <div className="text-heading-small text-text-primary">
                    #{rank}
                    </div>

                    <div className="text-caption-medium text-text-secondary uppercase">
                    Global Rank
                    </div>
                </div>
            )} */}
            <SmallButton variant="danger" onClick={onRemove}>
                Remove
            </SmallButton>
        </div>
    );
}
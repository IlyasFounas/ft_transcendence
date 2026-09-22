import UserProfile from "../user/UserProfile";
import type { UserProps } from "../../types/user.types";
import { X } from "lucide-react";

interface PendingFriendRowProps extends UserProps {
    onCancel?: () => void;
}

export default function PendingFriendRow({
    username,
    role,
    avatar,
    onCancel,
}: PendingFriendRowProps) {
  return (
        <div className="
            flex
            items-center
            justify-between
            bg-surface-base
            border
            border-border-default
            px-4
            xl:px-8
            h-16
            xl:h-20
            "
        >
            <UserProfile
                username={username}
                role={role}
                avatar={avatar}
            />

            <div className="flex items-center gap-6">
                <span className="
                    text-label-small
                    lg:text-heading-small
                    text-text-secondary
                    uppercase
                    "
                >
                    Pending
                </span>
                <button
                    onClick={onCancel}
                    aria-label="Cancel friend request"
                    className="
                        text-text-primary
                        cursor-pointer
                        hover:text-primary
                    "
                >
                    <X size={24} />
                </button>
            </div>
        </div>
    );
}
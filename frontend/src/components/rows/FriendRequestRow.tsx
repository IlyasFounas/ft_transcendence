import UserProfile from "../user/UserProfile";
import SmallButton from "../buttons/SmallButton";
import type { UserProps } from "../../types/user.types";
import { X } from "lucide-react";

interface FriendRequestRowProps extends UserProps {
    onAccept?: ()=> void;
    onReject?: ()=> void;
}

export default function FriendRequestRow({
    username,
    role,
    avatar,
    onAccept,
    onReject,
}: FriendRequestRowProps) {
    return (
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
            <div className="flex items-center gap-4">
                <SmallButton onClick={onAccept}>
                    Accept
                </SmallButton>
                <button
                    onClick={onReject}
                    aria-label="Reject friend request"
                    className="
                        text-text-primary
                        transition
                        duration-150
                        hover:drop-shadow-[0_0_6px_#35DFFF]
                    "
                >
                    <X size={24} className="cursor-pointer" />
                </button>
            </div>
        </div>
    );
}
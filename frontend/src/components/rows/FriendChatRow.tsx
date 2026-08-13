import type { UserProps, UserStatusType } from "../../types/user.types";
import Avatar from "../primitives/Avatar";

interface FriendChatRowProps extends UserProps{
    status: UserStatusType;
    lastMessage: string;
    timestamp: string;
    onOpenChat?: () => void;
}

export default function FriendChatRow ({
    username,
    role,
    avatar,
    status,
    lastMessage,
    timestamp,
    onOpenChat,
}: FriendChatRowProps) {
    const statusClasses = {
        online: "bg-status-green",
        "in-mission": "bg-status-yellow",
        offline: "bg-text-secondary",
    };

    return (
        <div className="
            grid
            grid-cols-[50px_minmax(0,1fr)_auto]
            gap-x-4
            items-center
            h-17.5
            bg-surface-base
            border
            border-surface-elevated
            px-2
            py-2
            cursor-pointer
            transition
            duration-150
            hover:bg-surface-elevated
            hover:shadow-[0_0_8px_rgba(53,223,255,0.15)]
            "
            onClick={onOpenChat}
        >
            <div className="relative">
                <Avatar avatar={avatar} size={50} />
                <span className={`absolute w-2 h-2 rounded-full bottom-2.5 right-0 ${statusClasses[status]}`} />
            </div>
            <div>
                <div className="text-body-primary text-text-primary uppercase">{username}</div>
                <div className="text-label-small text-primary uppercase">{role}</div>
                <div className="text-label-small text-text-secondary truncate">{lastMessage}</div>
            </div>
            <div className="text-label-small text-text-secondary self-start uppercase">
                {timestamp}
            </div>
        </div>
    );
}
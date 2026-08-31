import type { Friend } from "../../types/user.types";

import Avatar from "../primitives/Avatar";

interface CompactFriendRowProps {
    friend: Friend;
    onOpenChat: (friend: Friend) => void;
}

export default function CompactFriendRow({
    friend,
    onOpenChat,
}: CompactFriendRowProps) {
    
    const statusClasses = {
        online: "bg-status-green",
        "in-mission": "bg-status-yellow",
        offline: "bg-text-secondary",
    };

    return (
        <button
            type="button"
            onClick={() => onOpenChat(friend)}
            className="
                relative
                flex
                w-full
                justify-center
                py-2
                transition
                duration-150
                hover:bg-surface-elevated
            "
        >
            <div className="relative">
                <Avatar
                    avatar={friend.avatar}
                    size={50}
                />

                <span
                    className={`
                        absolute
                        bottom-1
                        right-0
                        h-2.5
                        w-2.5
                        rounded-full
                        ${statusClasses[friend.status]}
                    `}
                />
            </div>
        </button>
    );
}
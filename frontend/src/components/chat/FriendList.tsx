import FriendChatRow from "./FriendChatRow";
import type { Friend } from "../../types/user.types";
import CompactFriendRow from "./CompactFriendRow";

interface FriendsListProps {
    friends: Friend[];
    onOpenChat: (friend: Friend) => void;
}

export default function FriendsList({
    friends,
    onOpenChat,
}: FriendsListProps) {

    const onlineCount = friends.filter((friend) => friend.status === "online").length;

    return (
        <div
            className="
                flex
                h-full
                w-full
                flex-col
            "
        >

            <div className="shrink-0 px-4 py-3">
                <div className="text-heading-small text-text-primary uppercase">
                    Friends
                </div>

                <div className="text-body-primary text-status-green uppercase">
                    ● {onlineCount} Online
                </div>
            </div>

            <div
                className="
                    min-h-0
                    flex-1
                    mx-4
                    mb-4
                    overflow-y-auto
                    overflow-x-hidden
                    rounded-lg
                    border
                    border-border-default
                    custom-scrollbar
                "
            >
                {friends.map((friend, index) => (
                <div key={index} className="relative">

                    {/* Compact */}
                    <div className="block group-hover:hidden">
                        <CompactFriendRow
                            friend={friend}
                            onOpenChat={onOpenChat}
                        />
                    </div>

                    {/* Expanded */}
                    <div className="hidden group-hover:block">
                        <FriendChatRow
                            username={friend.username}
                            avatar={friend.avatar}
                            status={friend.status}
                            lastMessage={friend.lastMessage}
                            timestamp={friend.timestamp}
                            onOpenChat={() => onOpenChat(friend)}
                        />
                    </div>

                </div>
            ))}
            </div>
        </div>
    );
}``
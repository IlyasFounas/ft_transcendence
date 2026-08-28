import { useState } from "react";
import type { Friend } from "../../types/user.types";
import FriendsList from "./FriendList";
import Conversation from "./Conversation";

interface ChatPanelProps {
    friends: Friend[];
}

export default function ChatPanel({
    friends,
}: ChatPanelProps) {

    const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);

    return (
        <div className="h-full w-full">
            {selectedFriend === null ? (
                <FriendsList
                    friends={friends}
                    onOpenChat={setSelectedFriend}
                />
            ) : (
                <div
                    className="
                        flex
                        h-full
                        w-full
                        items-center
                        justify-center
                        border
                        border-border-default
                        bg-surface-base
                        text-text-primary
                    "
                >
                    <Conversation
                        friend={selectedFriend}
                        onBack={() => setSelectedFriend(null)}
                    />
                </div>
            )}
        </div>
    );
}
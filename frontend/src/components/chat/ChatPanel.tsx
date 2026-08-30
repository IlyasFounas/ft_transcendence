import { useState } from "react";
import type { Friend } from "../../types/user.types";
import FriendsList from "./FriendList";
import Conversation from "./Conversation";
import { conversations } from "../../data/conversation";
import type { Message } from "../../types/message.type";

interface ChatPanelProps {
    friends: Friend[];
}

const getCurrentTime = () => {
    return new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
    });
};

export default function ChatPanel({
    friends,
}: ChatPanelProps) {

    const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);
    
    const [conversationData, setConversationData] = useState(conversations);
    
    const [friendList, setFriendList] = useState(friends);

    const handleSendMessage = (message: string) => {
        if (!selectedFriend) {
            return;
        }

        const currentTime = getCurrentTime();

        const newMessage: Message = {
            id: Date.now(),
            variant: "sent",
            message,
            time: currentTime,
        };

        setConversationData((currentConversations) => ({
            ...currentConversations,
            [selectedFriend.username]: [
                ...(currentConversations[selectedFriend.username] ?? []),
                newMessage,
            ],
        }));

        setFriendList((currentFriends) => 
            currentFriends.map((friend) =>
                friend.username === selectedFriend.username
                    ? {
                        ...friend,
                        lastMessage: message,
                        timestamp: currentTime,
                    }
                    : friend
            )
        );
    };

    return (
        <div className="h-full w-full">
            {selectedFriend === null ? (
                <FriendsList
                    friends={friendList}
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
                        conversation={
                            conversationData[selectedFriend.username] ?? []
                        }
                        onBack={() => setSelectedFriend(null)}
                        onSendMessage={handleSendMessage}
                    />
                </div>
            )}
        </div>
    );
}
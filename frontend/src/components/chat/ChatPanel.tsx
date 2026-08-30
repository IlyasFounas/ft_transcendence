import { useState } from "react";
import type { Friend } from "../../types/user.types";
import FriendsList from "./FriendList";
import Conversation from "./Conversation";
import { conversations } from "../../data/conversation";

interface ChatPanelProps {
    friends: Friend[];
}

export default function ChatPanel({
    friends,
}: ChatPanelProps) {

    const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);
    
    const [conversationData, setConversationData] = useState(conversations);
    
    const handleSendMessage = (message: string) => {
        if (!selectedFriend) {
            return;
        }

        const newMessage = {
            id: Date.now(),
            variant: "sent" as const,
            message,
            time: "03:03 PM",
        };

        setConversationData((currentConversations) => ({
            ...currentConversations,
            [selectedFriend.username]: [
                ...(currentConversations[selectedFriend.username] ?? []),
                newMessage,
            ],
        }));
    };

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
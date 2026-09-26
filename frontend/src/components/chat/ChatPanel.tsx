import { useState, useEffect } from "react";
import type { Friend } from "../../types/user.types";
import FriendsList from "./FriendList";
import Conversation from "./Conversation";
import { conversations } from "../../data/conversation";
import type { Message } from "../../types/message.type";
import useChatWebSocket from "../../hooks/useChatWebSocket";
import { useAuth } from "../../context/AuthContext";

const getCurrentTime = () => {
    return new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
    });
};

export default function ChatPanel() {

    const { authFetch } = useAuth();
    const { chatFriends } = useChatWebSocket();
    console.log("ChatPanel chatFriends:", chatFriends);
    
    const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);
    const [conversationData, setConversationData] = useState(conversations);
    const [friendList, setFriendList] = useState<Friend[]>([]);
    
    useEffect(() => {
        const getFriends = async () => {
            const result = await authFetch("api/friend");

            if (!result) {
                return;
            }

            const formattedFriends: Friend[] = (result.data.friends ?? []).map(
                (friend: any) => ({
                    userId: friend.user_id,
                    username: friend.name,
                    avatar: friend.avatar || "/avatars/avatar1.png",
                })
            );

            setFriendList(formattedFriends);
        };

        getFriends();
    }, []);
    const friendsWithOnlineStatus = friendList.map(friend => {
        const chatFriend = chatFriends.find(
            chatFriend => chatFriend.user_id === friend.userId
        );
        
        return {
            ...friend,
            online: chatFriend?.online ?? false,
        };
    });

    console.log("Friends with WebSocket status:", friendsWithOnlineStatus);
    
    const handleSendMessage = (message: string) => {
        if (!selectedFriend) return;

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

        setFriendList((currentFriends) => {
        const updatedFriends = currentFriends.map((friend) =>
            friend.username === selectedFriend.username
            ? {
                ...friend,
                lastMessage: message,
                timestamp: currentTime,
                }
            : friend
        );

        const updatedFriend = updatedFriends.find(
            (friend) => friend.username === selectedFriend.username
        );

        const otherFriends = updatedFriends.filter(
            (friend) => friend.username !== selectedFriend.username
        );

        return updatedFriend
            ? [updatedFriend, ...otherFriends]
            : updatedFriends;
        });
    };

    return (
        <div className="flex h-full w-full flex-col">
            <div
                className={`
                min-h-0
                ${selectedFriend ? "hidden lg:block" : "block"}
                h-full
                lg:h-[45%]
                lg:shrink-0
                `}
            >
                <FriendsList
                    friends={friendsWithOnlineStatus}
                    onOpenChat={setSelectedFriend}
                />
            </div>

        <div
            className={`
            min-h-0
            ${selectedFriend ? "flex flex-1" : "hidden lg:flex lg:flex-1"}
            `}
        >
        {selectedFriend === null ? (
            <div className="flex h-full w-full items-center justify-center bg-surface-base text-text-primary">
                <div className="text-center">
                <div className="text-heading-small text-text-primary uppercase">
                    Select a conversation
                </div>
                <div className="mt-2 text-label-small text-text-secondary uppercase">
                    Choose an online agent to begin chatting.
                </div>
                </div>
            </div>
            ) : (
                <Conversation
                    friend={selectedFriend}
                    conversation={conversationData[selectedFriend.username] ?? []}
                    onBack={() => setSelectedFriend(null)}
                    onSendMessage={handleSendMessage}
                />
            )}
            </div>
        </div>
    );
}
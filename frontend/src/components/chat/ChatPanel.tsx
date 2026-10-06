import { useState, useEffect } from "react";
import type { Friend } from "../../types/user.types";
import FriendsList from "./FriendList";
import Conversation from "./Conversation";
import useChatWebSocket from "../../hooks/useChatWebSocket";
import { useAuth } from "../../context/AuthContext";

export default function ChatPanel() {

    const { authFetch } = useAuth();
    
    const { chatFriends, fetchMessages, sendMessage, messages, messageEnd, loadingOlderMessages } = useChatWebSocket();
    console.log("ChatPanel chatFriends:", chatFriends);
    
    const [selectedFriend, setSelectedFriend] = useState<Friend | null>(null);
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
            lastMessage: chatFriend?.last_message ?? "",
            timestamp: chatFriend?.date ?? "",
        };
    });

    useEffect(() => {
        if (!selectedFriend) {
            return;
        }

        const chatFriend = chatFriends.find(
            friend => friend.user_id === selectedFriend.userId
        );

        if (!chatFriend) {
            return;
        }

        setSelectedFriend(currentFriend => {
            if (!currentFriend) {
                return null;
            }
            
            return {
                ...currentFriend,
                online: chatFriend.online,
            };
        });
    }, [chatFriends]);

    console.log("Friends with WebSocket status:", friendsWithOnlineStatus);

    const handleSendMessage = (message: string) => {
        if (!selectedFriend) {
            return;
        }
        sendMessage(selectedFriend.userId, message);
    };

    const handleOpenChat = (friend: Friend) => {
        if (selectedFriend?.userId === friend.userId) {
            return;
        }
        setSelectedFriend(friend);
        fetchMessages(friend.userId);
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
                    onOpenChat={handleOpenChat}
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
                    conversation={messages}
                    messageEnd={messageEnd}
                    loadingOlderMessages={loadingOlderMessages}
                    fetchMessages={fetchMessages}
                    onBack={() => setSelectedFriend(null)}
                    onSendMessage={handleSendMessage}
                />
            )}
            </div>
        </div>
    );
}
import { useEffect, useRef, useState } from "react";
import type { Message } from "../types/message.type";
import { useAuth } from "../context/AuthContext";

interface ChatFriend {
    id: number;
    user_id: number;
    name: string;
    avatar: string | null;
    rank: string;
    online: boolean;
    last_message: string;
    date: string;
}

export default function useChatWebSocket() {
    
    const socketRef = useRef<WebSocket | null>(null);
    const loadingOlderRef = useRef(false);
    const socketTargetIdRef = useRef<number | null>(null);
    const { user } = useAuth();

    const [chatFriends, setChatFriends] = useState<ChatFriend[]>([]);
    const [messages, setMessages] = useState<Message[]>([]);
    const [messageEnd, setMessageEnd] = useState(0);
    const [loadingOlderMessages, setLoadingOlderMessages] = useState(false);
    const [messageError, setMessageError] = useState("");

    const fetchMessages = (targetId: number, end = 0) => {

        if (!socketRef.current) {
            return;
        }

        if (end === 0) {
            socketTargetIdRef.current = targetId;
        }

        if (end > 0) {
            if (loadingOlderRef.current) {
                return;
            }

            loadingOlderRef.current = true;
            setLoadingOlderMessages(true);
        }

        socketRef.current.send(
            JSON.stringify({
                action: "fetch_messages",
                payload: {
                    target_id: targetId,
                    end,
                },
            })
        );
    };

    const sendMessage = (targetId: number, message: string) => {
        if (!socketRef.current) {
            return ;
        }

        if ([...message].length > 1000) {
            setMessageError("Message is too long. Maximum 2000 characters.");
            return;
        }

        setMessageError("");

        socketRef.current.send(
            JSON.stringify({
                action: "send",
                payload: {
                    target_id: targetId,
                    message,
                }
            })
        );
    };

    useEffect(() => {
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";

        const socket = new WebSocket(`${protocol}//${window.location.host}/api/ws/chat`);

        socketRef.current = socket;

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            
            
            if (data.action === "friend_list") {
                if (Array.isArray(data.payload)) {

                    const payload = data.payload as ChatFriend[];
                    setChatFriends(payload);

                } else {
                    const payload = data.payload as ChatFriend[];

                    setChatFriends(currentFriends => {
                        const incomingFriend = payload as unknown as ChatFriend;
                        
                        const existingFriend = currentFriends.find(
                            friend => friend.user_id === incomingFriend.user_id
                        );

                        if (!existingFriend) {
                            return [...currentFriends, incomingFriend];
                        }

                        return currentFriends.map(friend =>
                            friend.user_id === incomingFriend.user_id
                            ? { ...friend, ...incomingFriend }
                            : friend
                        );
                    });
                }
            }

            if (data.action === "chats") {
                const chats = data.payload.chats ?? [];
                const end = data.payload.end ?? 0;

                setMessageEnd(end);

                if (loadingOlderRef.current) {
                    setMessages(currentMessages => [
                        ...[...chats].reverse(),
                        ...currentMessages,
                    ]);
                    loadingOlderRef.current = false;
                    setLoadingOlderMessages(false);
                } else {
                    setMessages([...chats].reverse());
                }
            }

            if (data.action === "new_message") {
                const newMessages = data.payload.chats ?? [];
                const latestMessage = newMessages[newMessages.length - 1];

                if (!latestMessage) {
                    return;
                }

                const friendId = 
                    latestMessage.user_id === user?.ID 
                        ? latestMessage.friend_id
                        : latestMessage.user_id;

                setChatFriends(currentFriends =>
                    currentFriends.map(friend =>
                        friend.user_id === friendId
                            ? {
                                ...friend,
                                last_message: latestMessage.message,
                                date: latestMessage.date,
                            }
                            : friend
                    )
                );

                if (socketTargetIdRef.current === friendId) {
                    setMessages(currentMessages => [
                        ...currentMessages,
                        ...newMessages,
                    ]);
                }
            }
        }

        return () => {
            socket.close();
            socketRef.current = null;
        };
    }, []);
    return {
        chatFriends,
        messages,
        fetchMessages,
        sendMessage,
        messageEnd,
        loadingOlderMessages,
        setMessageError,
        messageError,
    }
}
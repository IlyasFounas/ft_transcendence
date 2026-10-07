import { useEffect, useRef, useState } from "react";
import type { Message } from "../types/message.type";

interface ChatFriend {
    id: number;
    user_id: number;
    name: string;
    avatar: string | null;
    online: boolean;
    last_message: string;
    date: string;
}

export default function useChatWebSocket() {
    
    const socketRef = useRef<WebSocket | null>(null);
    const loadingOlderRef = useRef(false);

    const [chatFriends, setChatFriends] = useState<ChatFriend[]>([]);
    const [messages, setMessages] = useState<Message[]>([]);
    const [messageEnd, setMessageEnd] = useState(0);
    const [loadingOlderMessages, setLoadingOlderMessages] = useState(false);

    const fetchMessages = (targetId: number, end = 0) => {

        if (!socketRef.current) {
            return;
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
                    const payload = data.payload as ChatFriend;

                    setChatFriends(currentFriends =>
                        currentFriends.map(friend =>
                            friend.user_id === payload.user_id
                                ? { ...friend, online: payload.online }
                                : friend
                        )
                    );
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

                setMessages(currentMessages => [
                    ...currentMessages,
                    ...newMessages,
                ]);

                if (latestMessage) {
                    setChatFriends(currentFriends =>
                        currentFriends.map(friend =>
                            friend.user_id === latestMessage.user_id ||
                            friend.user_id === latestMessage.friend_id
                                ? {
                                    ...friend,
                                    last_message: latestMessage.message,
                                    date: latestMessage.date,
                                }
                                : friend
                        )
                    );
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
    }
}
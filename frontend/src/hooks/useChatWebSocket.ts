import { useEffect, useRef, useState } from "react";
import type { Message } from "../types/message.type";

interface ChatFriend {
    id: number;
    user_id: number;
    name: string;
    avatar: string | null;
    online: boolean;
}

export default function useChatWebSocket() {
    
    const socketRef = useRef<WebSocket | null>(null);

    const [chatFriends, setChatFriends] = useState<ChatFriend[]>([]);
    const [messages, setMessages] = useState<Message[]>([]);

    const fetchMessages = (targetId: number) => {
        console.log("fetchMessages called with:", targetId);

        if (!socketRef.current) {
            console.log("WebSocket is not connected");
            return;
        }

        console.log("Sending fetch_messages");

        socketRef.current.send(
            JSON.stringify({
                action: "fetch_messages",
                payload: {
                    target_id: targetId,
                    end: 0,
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

        socket.onopen = () => {
            console.log("Chat WebSocket connected");
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            
            if (data.title === "friend_list") {
                if (Array.isArray(data.payload)) {
                    
                    const payload = data.payload as ChatFriend[];
                    setChatFriends(payload);

                    console.log("Initial friend list:", payload);
                } else {
                    const payload = data.payload as ChatFriend;

                    setChatFriends(currentFriends =>
                        currentFriends.map(friend =>
                            friend.user_id === payload.user_id
                                ? { ...friend, online: payload.online }
                                : friend
                        )
                    );

                    console.log("Friend status update:", payload);
                }
            }

            if (data.title === "chats") {
                console.log("This is a chat history:", data.payload);
                setMessages(data.payload.chats ?? []);
            }

            if (data.title === "new_message") {
                console.log("This is a new message:", data.payload);
                setMessages(currentMessages => [
                    ...currentMessages,
                    ...(data.payload.chats ?? []),
                ]);
            }
        }
        
        socket.onerror = (error) => {
            console.error("Chat WebSocket error:", error);
        };

        return () => {
            console.log("WebSocket cleanup");
            socket.close();
            socketRef.current = null;
        };
    }, []);
    return {
        chatFriends,
        messages,
        fetchMessages,
        sendMessage,
    }
}
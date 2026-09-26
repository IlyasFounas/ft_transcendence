import { useEffect, useRef, useState } from "react";

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
            }

            if (data.title === "new_message") {
                console.log("This is a new message:", data.payload);
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
    }
}
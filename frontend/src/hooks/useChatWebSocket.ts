import { useEffect, useRef } from "react";

export default function useChatWebSocket() {
    
    const socketRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";

        const socket = new WebSocket(`${protocol}//${window.location.host}/api/ws/chat`);

        socketRef.current = socket;

        socket.onopen = () => {
            console.log("Chat WebSocket connected");
        };

        socket.onerror = (error) => {
            console.error("Chat WebSocket error:", error);
        };

        socket.onclose = () => {
            socket.close();
            socketRef.current = null;
        };
    }, []);
}
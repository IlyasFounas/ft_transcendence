import { Outlet } from "react-router-dom";
import { ChatProvider } from "../context/ChatContext";

export function ProtectedChatLayout() {
    return (
        <ChatProvider>
            <Outlet />
        </ChatProvider>
    );
}
import { createContext, useContext, type ReactNode } from "react";
import useChatWebSocket from "../hooks/useChatWebSocket";

type ChatContextType = ReturnType<typeof useChatWebSocket>;

const ChatContext = createContext<ChatContextType | null>(null);

export function ChatProvider({ children }: { children: ReactNode }) {
	const chat = useChatWebSocket();

	return (
		<ChatContext.Provider value={chat}>
			{children}
		</ChatContext.Provider>
	);
}

export function useChat() {
	const context = useContext(ChatContext);

	if (!context) {
		throw new Error("useChat must be used within ChatProvider");
	}

	return context;
}

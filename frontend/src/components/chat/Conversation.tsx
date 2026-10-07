import { ArrowLeft } from "lucide-react";
import ChatBoxCard from "../cards/ChatBox";
import type { Friend } from "../../types/user.types";
import Avatar from "../primitives/Avatar";
import { useState, useLayoutEffect, useRef, } from "react";
import MessageInput from "./MessageInput";
import type { Message } from "../../types/message.type";
import { useAuth } from "../../context/AuthContext";

interface ConversationProps {
    friend: Friend;
    conversation: Message[];
    messageEnd: number;
    loadingOlderMessages: boolean;
    fetchMessages: (targetId: number, end?: number) => void;
    onBack: () => void;
    onSendMessage: (message: string) => void;
}

export default function Conversation({
    friend,
    conversation,
    messageEnd,
    loadingOlderMessages,
    fetchMessages,
    onBack,
    onSendMessage,
}: ConversationProps) {

    const { user } = useAuth();

    const statusClasses = {
        online: "bg-status-green",
        "in-mission": "bg-status-yellow",
        offline: "bg-text-secondary",
    };

    const [message, setMessage] = useState("");

    const handleSend = () => {
        if (message.trim() === "") {
            return;
        }

        onSendMessage(message);
        setMessage("");
    };

    const messagesEndRef = useRef<HTMLDivElement>(null);
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    const prevMessageCountRef = useRef(0);
    const prevScrollHeightRef = useRef(0);

    useLayoutEffect(() => {
        const container = messagesContainerRef.current;

        if (!container) {
            return;
        }

        const prevCount = prevMessageCountRef.current;

        if (prevCount === 0 && conversation.length > 0) {
            messagesEndRef.current?.scrollIntoView({
                behavior: "auto",
                block: "end",
            });
        }

        if (conversation.length > prevCount) {
            const heightDif = container.scrollHeight - prevScrollHeightRef.current;

            if (prevScrollHeightRef.current > 0) {
                container.scrollTop += heightDif;
            }
        }
        
        prevMessageCountRef.current = conversation.length;
        prevScrollHeightRef.current = container.scrollHeight;
    }, [conversation]);

    return (
        <div
            className="
                flex
                h-full
                w-full
                min-w-0
                flex-col
                overflow-hidden
            "
        >
            <div className="shrink-0 p-2">
                <div className="flex items-center gap-3">
                    <button
                        type="button"
                        onClick={onBack}
                        className="
                            text-text-secondary
                            transition
                            duration-150
                            hover:text-primary
                        "
                    >
                        <ArrowLeft size={20} />
                    </button>
                    
                    <div className="relative shrink-0">
                        <Avatar avatar={friend.avatar} size={50} />

                        <span
                            className={`
                                absolute
                                bottom-2
                                right-0
                                h-2.5
                                w-2.5
                                rounded-full
                                ${friend.online ? statusClasses.online : statusClasses.offline}
                            `}
                        />
                    </div>

                    <div className="min-w-0">
                        <div className="text-heading-small text-text-primary uppercase truncate">
                            {friend.username}
                        </div>
                        <div className="text-label-small text-primary uppercase">
                            {friend.role}
                        </div>
                    </div>
                </div>
            </div>

            <div
                ref={messagesContainerRef}
                className="
                    min-h-0
                    min-w-0
                    flex-1
                    overflow-y-auto
                    overflow-x-hidden
                    rounded-lg
                    border
                    border-border-default
                    px-2
                    py-4
                    mx-4
                    mt-0
                    mb-4
                    flex
                    flex-col
                    gap-4
                    custom-scrollbar
                "
                onScroll={(event) => {
                    if (event.currentTarget.scrollTop === 0 && 
                        messageEnd > 0 && !loadingOlderMessages
                    ) {
                        prevScrollHeightRef.current = event.currentTarget.scrollHeight;

                        console.log("Loading older message");
                        fetchMessages(friend.userId, messageEnd);
                    }
                }}
            >
                {conversation.map((message) => (
                    <ChatBoxCard
                        key={message.id}
                        variant={message.user_id === user?.ID ? "sent" : "received"}
                        message={message.message}
                        time={new Date(message.date).toLocaleTimeString([], {
                            hour: "2-digit",
                            minute: "2-digit",
                        })}
                    />
                ))}

                <div ref={messagesEndRef} />
            </div>

            <div className="shrink-0">
                <MessageInput
                    value={message}
                    onChange={setMessage}
                    onSend={handleSend}
                />
            </div>
            
        </div>
    );
}
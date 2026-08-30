import { ArrowLeft } from "lucide-react";
import ChatBoxCard from "../cards/ChatBox";
import type { Friend } from "../../types/user.types";
import Avatar from "../primitives/Avatar";
import { useState, useEffect, useRef, } from "react";
import MessageInput from "./MessageInput";
import type { Message } from "../../types/message.type";

interface ConversationProps {
    friend: Friend;
    conversation: Message[];
    onBack: () => void;
    onSendMessage: (message: string) => void;
}

export default function Conversation({
    friend,
    conversation,
    onBack,
    onSendMessage,
}: ConversationProps) {

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

    // auto-scroll to the newest message
    const messagesEndRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: "smooth",
        });
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
                bg-surface-base
                border
                border-border-default
            "
        >
            <div className="shrink-0 border-b border-border-default p-4">
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
                                ${statusClasses[friend.status]}
                            `}
                        />
                    </div>

                    <div>
                        <div className="text-heading-small text-text-primary uppercase">
                            {friend.username}
                        </div>
                        <div className="text-label-small text-primary uppercase">
                            {friend.role}
                        </div>
                    </div>
                </div>
            </div>

            <div
                className="
                    min-h-0
                    min-w-0
                    flex-1
                    overflow-y-auto
                    overflow-x-hidden
                    py-4
                    flex
                    flex-col
                    gap-2
                    custom-scrollbar
                "
            >
                {conversation.map((message) => (
                    <ChatBoxCard
                        key={message.id}
                        variant={message.variant}
                        message={message.message}
                        time={message.time}
                    />
                ))}

                <div ref={messagesEndRef} />
            </div>

            <MessageInput
                value={message}
                onChange={setMessage}
                onSend={handleSend}
            />
            
        </div>
    );
}
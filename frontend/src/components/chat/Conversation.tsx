import { ArrowLeft } from "lucide-react";
import ChatBoxCard from "../cards/ChatBox";
import type { Friend } from "../../types/user.types";
import Avatar from "../primitives/Avatar";

interface ConversationProps {
    friend: Friend;
    onBack: () => void;
}

export default function Conversation({
    friend,
    onBack,
}: ConversationProps) {
    const statusClasses = {
        online: "bg-status-green",
        "in-mission": "bg-status-yellow",
        offline: "bg-text-secondary",
    };

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
                <ChatBoxCard
                    variant="received"
                    message="Hello! How are you doing?"
                    time="03:01 PM"
                />
                <ChatBoxCard
                    variant="sent"
                    message="I'm good, how are you?"
                    time="03:02 PM"
                />
                <ChatBoxCard
                    variant="received"
                    message="Great! Do you wanna join game now?"
                    time="03:02 PM"
                />
            </div>
        </div>
    );
}
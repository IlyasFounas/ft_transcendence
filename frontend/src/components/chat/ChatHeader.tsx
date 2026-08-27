import { ArrowLeft, X } from "lucide-react";
import type { UserProps } from "../../types/user.types";
import type { UserStatusType } from "../../types/user.types";

interface ChatHeaderProp extends UserProps {
    status: UserStatusType;
    onBack: () => void;
    onClose: () => void;
}

export function ChatHeader({
    avatar,
    username,
    role,
    status,
    onBack,
    onClose,
}: ChatHeaderProp) {
    return (
        <div className="border-b border-border-default pb-3">
            <div className="flex items-center justify-between">
                <button
                    type="button"
                    onClick={onBack}
                    className="
                        flex
                        items-center
                        gap-2
                        text-text-secondary
                        transition
                        duration-150
                        hover:text-primary
                    "
                >
                    <ArrowLeft size={20} />
                    <span className="text-body-primary uppercase">
                        Friends
                    </span>
                </button>

                <button
                    type="button"
                    onClick={onClose}
                    className="
                        text-text-secondary
                        transition
                        duration-150
                        hover:text-primary
                    "
                >
                    <X size={20} />
                </button>
            </div>

            <div className="mt-3 flex items-center gap-3">
                <div className="relative">
                    <img
                        src={avatar}
                        alt={username}
                        className="
                            h-10
                            w-10
                            rounded-full
                            object-cover
                        "
                    />
                    <span
                        className={`
                            absolute
                            bottom-0
                            right-0
                            h-2.5
                            w-2.5
                            rounded-full
                            border
                            border-surface-base
                            ${
                                status === "online"
                                    ? "bg-status-green"
                                    : status === "in-mission"
                                    ? "bg-status-yellow"
                                    : "bg-text-secondary"
                            }
                        `}
                    />
                </div>
                <div>
                    <div className="text-body-primary text-text-primary uppercase">
                        {username}
                    </div>
                    <div className="text-[9px] text-primary uppercase">
                        {role}
                    </div>

                </div>

            </div>

        </div>
    );
}
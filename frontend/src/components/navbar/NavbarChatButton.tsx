import { Users, X } from "lucide-react";

interface NavbarChatButtonProps {
    isChatOpen: boolean;
    onClick: () => void;
}

export default function NavbarChatButton({
    isChatOpen,
    onClick,
}: NavbarChatButtonProps) {
    return (
        <button
            type="button"
            onClick={onClick}
            className="
                flex
                h-12
                w-12
                items-center
                justify-center
                text-text-secondary
                transition
                duration-150
                hover:text-primary
            "
        >
            {isChatOpen ? (
                <X size={24} />
            ) : (
                <Users size={24} />
            )}
        </button>
    );
}
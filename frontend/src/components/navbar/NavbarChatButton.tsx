import { Users } from "lucide-react";

interface NavbarChatButtonProps {
    onClick: () => void;
}

export default function NavbarChatButton({
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
            <Users size={24} />
        </button>
    );
}
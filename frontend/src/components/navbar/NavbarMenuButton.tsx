import { Menu } from "lucide-react";

interface NavbarMenuButtonProps {
    onClick: () => void;
}

export default function NavbarMenuButton({
    onClick,
}: NavbarMenuButtonProps) {
    return (
        <button
            type="button"
            onClick={onClick}
            className="
                flex
                items-center
                justify-center
                text-text-secondary
                transition
                duration-150
                hover:text-primary
                hover:drop-shadow-[0_0_6px_rgba(53,223,255,0.5)]
            "
        >
            <Menu size={28} />
        </button>
    );
}
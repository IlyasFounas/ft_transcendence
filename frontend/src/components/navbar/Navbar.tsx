import { useState } from "react";
import { NavLink } from "react-router-dom";
import type { UserStatusType } from "../../types/user.types";
import { useAuth } from "../../context/AuthContext";

import NavbarUser from "./NavbarUser";
import NavbarMenuButton from "./NavbarMenuButton";
import NavbarMobileMenu from "./NavbarMobileMenu";
import NavbarChatButton from "./NavbarChatButton";

interface NavbarProps {
    status: UserStatusType;
    isChatOpen?: boolean;
    onToggleChat?: () => void;
    showChat: boolean;
}

export function Navbar({
    status,
    isChatOpen,
    onToggleChat,
    showChat,
}: NavbarProps) {

    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const { user } = useAuth();

    return (
        <nav
            className="
                w-full
                h-20
                lg:h-25
                shrink-0
                bg-surface-base
                border
                border-border-default
            "
        >
            {/* Navbar for desktop */}
            <div className="
                hidden
                xl:flex
                h-full
                items-center
                px-8
            ">
                {/* logo */}
                <div className="w-50 lg:70 xl:w-80">
                    <NavLink to="/home">
                        <img
                            src="/logos/atlas-logo.png"
                            alt="Atlas"
                            className="h-15 lg:h-20"
                        />
                    </NavLink>
                </div>
                
                {/* menu */}
                <div className="
                    flex-1
                    flex
                    items-center
                    justify-center
                    gap-8
                    xl:gap-20
                    text-heading-small
                    lg:text-heading-medium
                    text-text-primary
                    uppercase
                ">
                    <NavLink
                        to="/home"
                        className={({ isActive }) => `
                            relative
                            h-full
                            flex
                            items-center
                            whitespace-nowrap
                            transition
                            duration-150
                            hover:text-primary
                            ${isActive ? "nav-active-line" : ""}
                        `}
                    >
                        Command Center
                    </NavLink>

                    <NavLink
                        to="/multiplayer"
                        className={({ isActive }) => `
                            relative
                            h-full
                            flex
                            items-center
                            whitespace-nowrap
                            transition
                            duration-150
                            hover:text-primary
                            ${isActive ? "nav-active-line" : ""}
                        `}
                    >
                        Multiplayer
                    </NavLink>

                    <NavLink
                        to="/leaderboard"
                        className={({ isActive }) => `
                            relative
                            h-full
                            flex
                            items-center
                            whitespace-nowrap
                            transition
                            duration-150
                            hover:text-primary
                            ${isActive ? "nav-active-line" : ""}
                        `}
                    >
                        Leaderboard
                    </NavLink>
                </div>

                {/* profile */}
                <div className="w-50 lg:70 xl:w-80 flex justify-end">
                    <NavbarUser
                        username={`Agent ${user?.Username}`}
                        status={status}
                        avatar={user?.Avatar ?? "/avatars/avatar1.png"}
                    />
                </div>
            </div>

            {/* Navbar for tablet and phone */}
            <div className="
                flex
                xl:hidden
                h-full
                items-center
                px-6
                "
            >
                {/* menu */}
                <div className="w-12">
                    <NavbarMenuButton
                        onClick={() => setIsMenuOpen(!isMenuOpen)}
                    />
                </div>

                {/* logo */}
                <div className="flex-1 flex justify-center">
                    <NavLink to="/home">
                        <img
                            src="/logos/atlas-logo.png"
                            alt="Atlas"
                            className="h-12"
                        />
                    </NavLink>
                </div>

                {/* chat */}
                <div className="w-12 flex justify-end">
                    {showChat && (
                        <NavbarChatButton
                            isChatOpen={isChatOpen ?? false}
                            onClick={onToggleChat ?? (() => {})}
                        />
                    )}
                </div>
            </div>
            
            {/* open menu */}
            <NavbarMobileMenu
                isOpen={isMenuOpen}
                onClose={() => setIsMenuOpen(false)}
            />
        </nav>
    );
}
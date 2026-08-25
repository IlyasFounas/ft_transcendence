import { NavLink } from "react-router-dom";
import NavbarUser from "./NavbarUser";
import type { UserStatusType } from "../../types/user.types";
import { useState } from "react";
import NavbarMenuButton from "./NavbarMenuButton";
import NavbarMobileMenu from "./NavbarMobileMenu";

interface NavbarProps {
    username: string;
    avatar: string;
    status: UserStatusType;
}

export function Navbar({
    username,
    avatar,
    status,
}: NavbarProps) {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    return (
        <nav
            className="
                w-full
                h-20
                bg-surface-base
                border
                border-border-default
            "
        >
            {/* Navbar for desktop */}
            <div className="
                hidden
                lg:flex
                h-full
                items-center
                px-8
            ">
                <div className="w-75">
                    <NavLink to="/home">
                        <img
                            src="/logos/atlas-logo.png"
                            alt="Atlas"
                            className="h-15"
                        />
                    </NavLink>
                </div>

                <div className="
                    flex-1
                    flex
                    items-center
                    justify-center
                    gap-20
                    text-heading-small
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
                <div className="w-75 flex justify-end">
                    <NavbarUser
                        username={username}
                        status={status}
                        avatar={avatar}
                    />
                </div>
            </div>

            {/* Navbar for tablet and phone */}
            <div className="
                flex
                lg:hidden
                h-full
                items-center
                px-6
                "
            >
                <div className="w-12">
                    <NavbarMenuButton
                        onClick={() => setIsMenuOpen(!isMenuOpen)}
                    />
                </div>
                <div className="flex-1 flex justify-center">
                    <NavLink to="/home">
                        <img
                            src="/logos/atlas-logo.png"
                            alt="Atlas"
                            className="h-12"
                        />
                    </NavLink>
                </div>

                <div className="w-12 flex justify-end">
                    {/* chat */}
                </div>
            </div>
            <NavbarMobileMenu
                isOpen={isMenuOpen}
                onClose={() => setIsMenuOpen(false)}
            />
        </nav>
    );
}
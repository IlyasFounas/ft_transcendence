import { NavLink } from "react-router-dom";
import { House, Users, Trophy, User, LogOut, } from "lucide-react";

interface NavbarMobileMenuProps {
    isOpen: boolean;
    onClose: () => void;
}

export default function NavbarMobileMenu({
    isOpen,
    onClose,
}: NavbarMobileMenuProps) {
    return (
        <>
            { isOpen && (
                <div
                    className="
                        fixed
                        inset-0
                        z-30
                        bg-black/20
                        lg:hidden
                    "
                    onClick={onClose}
                />
            )}
            
            <div
                className={`
                    fixed
                    left-0
                    top-20
                    bottom-0
                    z-50
                    lg:hidden

                    w-full
                    md:w-72

                    bg-[linear-gradient(90deg,#133867_0%,#020A14_100%)]

                    transition-transform
                    duration-300
                    ease-out

                    ${isOpen
                        ? "translate-x-0"
                        : "-translate-x-full"
                    }
                `}
            >
                <div className="
                    flex
                    flex-col
                    py-4
                ">
                    <NavLink
                        to="/home"
                        onClick={onClose}
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-3
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-primary
                        "
                    >
                        <House size={16} />
                        <span>Command Center</span>
                    </NavLink>

                    <NavLink
                        to="/multiplayer"
                        onClick={onClose}
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-3
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-primary
                        "
                    >
                        <Users size={16} />
                        <span>Multiplayer</span>
                    </NavLink>

                    <NavLink
                        to="/leaderboard"
                        onClick={onClose}
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-3
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-primary
                        "
                    >
                        <Trophy size={16} />
                        <span>Leaderboard</span>
                    </NavLink>

                    <NavLink
                        to="/profile"
                        onClick={onClose}
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-3
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-primary
                        "
                    >
                        <User size={16} />
                        <span>Profile</span>
                    </NavLink>

                    <button
                        type="button"
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-3
                            text-body-primary
                            text-status-red
                            uppercase
                            hover:brightness-125
                            cursor-pointer
                        "
                    >
                        <LogOut size={16} />
                        <span>Sign Out</span>
                    </button>
                </div>
            </div>
        </>
    );
}
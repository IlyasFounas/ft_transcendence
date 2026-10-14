import { NavLink } from "react-router-dom";
import { House, Users, Trophy, User, LogOut, } from "lucide-react";
import { useAuth } from "../../context/AuthContext";

interface NavbarMobileMenuProps {
    isOpen: boolean;
    onClose: () => void;
}

export default function NavbarMobileMenu({
    isOpen,
    onClose,
}: NavbarMobileMenuProps) {
    
    const { logout } = useAuth();

    return (
        <>
            {/* add background overlay */}
            { isOpen && (
                <div
                    className="
                        fixed
                        inset-0
                        z-30
                        bg-black/20
                        xl:hidden
                    "
                    onClick={onClose}
                />
            )}
            
            {/* sliding menu */}
            <div
                className={`
                    fixed
                    left-0
                    top-20
                    bottom-0
                    z-50
                    xl:hidden
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

                {/* menu items */}
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
                        onClick={logout}
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
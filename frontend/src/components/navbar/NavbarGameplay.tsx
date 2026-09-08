import { useState } from "react";
import { NavLink } from "react-router-dom";
import NavbarMenuButton from "./NavbarMenuButton";
import NavbarMobileMenu from "./NavbarMobileMenu";

export default function NavbarGameplay() {

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
                {/* logo */}
                <div className="w-75">
                        <img
                            src="/logos/atlas-logo.png"
                            alt="Atlas"
                            className="h-15"
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
            </div>
            
            {/* open menu */}
            <NavbarMobileMenu
                isOpen={isMenuOpen}
                onClose={() => setIsMenuOpen(false)}
            />
        </nav>
    );
}
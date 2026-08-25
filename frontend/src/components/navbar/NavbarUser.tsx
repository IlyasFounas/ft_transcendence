import { ChevronDown, User, LogOut, } from "lucide-react";
import { NavLink } from "react-router-dom";
import Avatar from "../primitives/Avatar";

interface NavbarUserProps {
    username: string;
    avatar: string;
    status: string;
}

export default function NavbarUser({
    username,
    status,
    avatar,
}: NavbarUserProps) {
    return (
        <div className="
            group
            relative
            flex
            items-center
            gap-3
        ">
            <NavLink
                to="/profile"
                className="
                    flex
                    items-center
                    gap-3
                    transition
                    duration-150
                    hover:brightness-125
                "
            >
                <Avatar
                    avatar={avatar}
                    size={45}
                />

                <div>
                    <div className="
                        text-heading-small
                        text-text-primary
                        uppercase
                    ">
                        {username}
                    </div>

                    <div className="
                        text-body-primary
                        text-status-green
                        uppercase
                    ">
                        ● {status}
                    </div>
                </div>
            </NavLink>

            <button
                type="button"
                className="
                    p-1
                    text-text-secondary
                    transition
                    duration-150
                    group-hover:text-primary
                "
            >
                <ChevronDown
                    size={18}
                    className="
                        transition-transform
                        duration-150
                        group-hover:rotate-180
                    "
                />
            </button>

            <div className="
                absolute
                right-0
                top-full
                hidden
                pt-3
                group-hover:block
            ">
                <div className="
                    w-60
                    rounded-lg
                    border
                    border-border-default
                    bg-surface-elevated
                    shadow-[0_0_15px_rgba(0,183,255,0.2)]
                    overflow-hidden
                ">
                    <NavLink
                        to="/profile"
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-4
                            text-body-primary
                            text-text-primary
                            uppercase
                            transition
                            duration-150
                            hover:text-primary
                        "
                    >
                        <User size={18} />
                        <span>profile</span>
                    </NavLink>

                    <button
                        type="button"
                        className="
                            flex
                            items-center
                            gap-3
                            px-5
                            py-4
                            text-body-primary
                            text-status-red
                            uppercase
                            transition
                            duration-150
                            hover:brightness-125
                            cursor-pointer
                        "
                    >
                        <LogOut size={18} />
                        <span>log out</span>
                    </button>
                </div>
            </div>
        </div>
    );
}
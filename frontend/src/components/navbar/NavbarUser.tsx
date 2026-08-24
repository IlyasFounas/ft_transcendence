import { ChevronDown } from "lucide-react";
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
                    w-40
                    rounded-lg
                    border
                    border-border-default
                    bg-surface-elevated
                    shadow-[0_0_15px_rgba(0,183,255,0.2)]
                ">
                    <NavLink
                        to="/profile"
                        className="
                            block
                            px-4
                            py-3
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-primary
                        "
                    >
                        profile
                    </NavLink>

                    <button
                        type="button"
                        className="
                            w-full
                            px-4
                            py-3
                            text-left
                            text-body-primary
                            text-text-primary
                            uppercase
                            hover:text-status-red
                        "
                    >
                        log out
                    </button>
                </div>
            </div>
        </div>
    );
}
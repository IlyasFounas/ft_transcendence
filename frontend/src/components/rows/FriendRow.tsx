import UserProfile from "../user/UserProfile";
import type { UserProps } from "../../types/user.types";

export default function FriendRow({
    username,
    role,
    avatar,
}: UserProps) {
    return (
        <div className="
            flex
            items-center
            justify-between
            bg-surface-base
            border
            border-border-default
            px-4
            h-16
            "
        >
            <UserProfile
                username={username}
                role={role}
                avatar={avatar}
            />

            <div className="flex items-center gap-6">
                <span className="
                    text-heading-small
                    text-text-secondary
                    uppercase
                    "
                >
                    Friends
                </span>
            </div>
        </div>
    );
}
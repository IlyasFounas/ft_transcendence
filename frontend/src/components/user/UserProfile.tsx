import Avatar from "../primitives/Avatar";
import type { UserProps } from "../../types/user.types";

export default function UserProfile({
    username,
    role,
    avatar,
}: UserProps) {
    return (
        <div className="flex items-center gap-4 min-w-0 flex-1">
            <Avatar avatar={avatar} size={45}/>
            <div className="min-w-0 flex-1">
                <div className="text-body-primary md:text-heading-small xl:text-heading-medium text-text-primary uppercase truncate">
                    <span>Agent </span>
                    {username}
                </div>
                <div className="text-label-small md:text-body-primary xl:text-heading-small text-text-secondary uppercase">{role}</div>
            </div>
        </div>
    );
}
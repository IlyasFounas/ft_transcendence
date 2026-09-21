import Avatar from "../primitives/Avatar";
import type { UserProps } from "../../types/user.types";

export default function UserProfile({
    username,
    role,
    avatar,
}: UserProps) {
    return (
        <div className="flex items-center gap-4">
            <Avatar avatar={avatar} size={45}/>
            <div>
                <div className="text-heading-small text-text-primary uppercase">
                    <span>Agent </span>
                    {username}
                </div>
                <div className="text-body-primary text-text-secondary uppercase">{role}</div>
            </div>
        </div>
    );
}
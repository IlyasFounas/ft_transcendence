import Avatar from "../primitives/Avatar";

interface UserProfileProps {
    username: string;
    role: string;
    avatar: string;
}

export default function UserProfile({
    username,
    role,
    avatar,
}: UserProfileProps) {
    return (
        <div className="flex items-center gap-4">
            <Avatar avatar={avatar} size={45}/>
            <div>
                <div className="text-heading-small text-text-primary uppercase">{username}</div>
                <div className="text-body-primary text-text-secondary uppercase">{role}</div>
            </div>
        </div>
    );
}
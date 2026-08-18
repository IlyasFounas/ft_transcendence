import { Mail } from "lucide-react"
import PrimaryButton from "../buttons/PrimaryButton";

interface AuthenticationWelcomeProps {
    onGoogleLogin?: () => void;
    onEmailSignUp?: () => void;
    onEmailLogin?: () => void;
}

export default function AuthenticationWelcome ({
    onGoogleLogin,
    onEmailSignUp,
    onEmailLogin,
}: AuthenticationWelcomeProps) {
    return (
        <div className="
            flex
            flex-col
            w-full
            gap-[50px]
            "
        >
            <div className="flex flex-col gap-[30px]">
                <button className="
                    flex
                    h-[62px]
                    items-center
                    justify-center
                    gap-2
                    rounded-2xl
                    bg-white
                    text-black
                    text-heading-medium
                    transition
                    duration-200
                    hover:shadow-[0_0_12px_rgba(255,255,255,0.35)]
                    hover:scale-[1.02]
                    active:scale-[0.98]
                    cursor-pointer
                    "
                    onClick={onGoogleLogin}
                >
                    <img
                        src="icons/gmail.svg"
                        alt=""
                        className="w-5 h-5"
                    />
                    <span>Continue with Google</span>
                </button>
                <button className="
                    flex
                    h-[62px]
                    items-center
                    justify-center
                    gap-2
                    rounded-2xl
                    bg-primary
                    text-black
                    text-heading-medium
                    transition
                    duration-200
                    hover:shadow-[0_0_12px_rgba(53,223,255,0.35)]
                    hover:scale-[1.02]
                    active:scale-[0.98]
                    cursor-pointer
                    "
                    onClick={onEmailSignUp}
                >
                    <Mail size={20} />
                    <span>Continue with Email</span>
                </button>
            </div>
            <div className="flex gap-[16px] items-center">
                <hr className="w-full border-t border-text-secondary" />
                <span className="text-body-primary text-text-secondary uppercase">
                    or
                </span>
                <hr className="w-full border-t border-text-secondary" />
            </div>
            <div className="flex items-center justify-between">
                <div className="text-heading-small text-text-primary">
                    Already have an account ?
                </div>
                <PrimaryButton size="small" onClick={onEmailLogin}>
                    Login
                </PrimaryButton>
            </div>
        </div>
    );
}
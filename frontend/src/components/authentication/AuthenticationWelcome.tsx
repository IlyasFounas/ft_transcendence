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
            gap-5
            md:py-12
            "
        >
            <div className="flex flex-col gap-5 md:px-6">
                <button className="
                    flex
                    items-center
                    justify-center
                    gap-2
                    rounded-2xl
                    bg-white
                    text-black
                    transition
                    duration-200
                    hover:shadow-[0_0_12px_rgba(255,255,255,0.35)]
                    hover:scale-[1.02]
                    active:scale-[0.98]
                    cursor-pointer
                    
                    h-13
                    text-body-primary
                    md:text-heading-medium
                    md:h-15
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
                    items-center
                    justify-center
                    gap-2
                    rounded-2xl
                    bg-primary
                    text-black
                    transition
                    duration-200
                    hover:shadow-[0_0_12px_rgba(53,223,255,0.35)]
                    hover:scale-[1.02]
                    active:scale-[0.98]
                    cursor-pointer

                    h-13
                    text-body-primary
                    md:text-heading-medium
                    md:h-15
                    "
                    onClick={onEmailSignUp}
                >
                    <Mail size={20} />
                    <span>Continue with Email</span>
                </button>
            </div>
            <div className="flex w-full gap-4 items-center sm:px-6">
                <hr className="h-px flex-1 border-text-secondary" />
                <span className="text-body-primary text-text-secondary uppercase">
                    or
                </span>
                <hr className="h-px flex-1 border-text-secondary" />
            </div>
            <div className="flex flex-col items-center px-2">
                <div className="uppercase text-label-small text-text-primary md:text-body-primary">
                    Already have an account ?
                </div>
                <PrimaryButton size="small" onClick={onEmailLogin}>
                    log in
                </PrimaryButton>
            </div>
        </div>
    );
}
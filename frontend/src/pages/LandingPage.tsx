import { useNavigate } from "react-router-dom";
import PrimaryButton from "../components/buttons/PrimaryButton";

export default function LandingPage() {
    const navigate = useNavigate();

    return (
        <main
            className="
                min-h-dvh
                w-full
                bg-surface-base
                bg-cover
                bg-center
                bg-[url('/images/landing-bg.png')]
                bg-no-repeat
                flex
                items-center
                justify-center
                px-5
                py-10
            "
        >
            <div
                className="
                    flex
                    w-full
                    flex-col
                    items-center
                    text-center
                    max-w-sm
                    sm:max-w-md
                    lg:max-w-xl
                "
            >
                <div className="flex flex-col items-center">
                    <img
                        src="/logos/atlas-logo-landing.png"
                        alt="ATLAS"
                        className="w-60 md:w-93 lg:w-146"
                    />
                </div>
                <div 
                    className="
                        flex
                        items-center
                        justify-center
                        uppercase
                        text-white
                        py-4

                        gap-5
                        text-heading-small
                        
                        md:text-heading-medium
                        md:gap-8

                        lg:text-heading-large
                    "
                >
                    <span>Analyze.</span>
                    <span>Locate.</span>
                    <span>Report.</span>
                </div>

                <div
                    className="
                        flex
                        w-full
                        flex-col
                        items-center
                        mt-2
                        
                        max-w-xs
                        lg:max-w-md
                    "
                >
                    <PrimaryButton
                        size="medium"
                        onClick={() => navigate("/login")}
                    >
                        Login
                    </PrimaryButton>

                    <PrimaryButton 
                        size="small"
                        onClick={() => navigate("/game")}
                    >
                      Start Guest Mode
                    </PrimaryButton>
                </div>

                <button
                    type="button"
                    className="
                        mt-7
                        md:mt-7
                        cursor-pointer
                        text-body-primary
                        md:text-heading-small
                        lg:text-heading-medium
                        uppercase
                        text-text-primary
                        transition
                        hover:text-primary
                    "
                >
                    How to play <span>▶</span>
                </button>
            </div>
        </main>
    );
}
import { useNavigate } from "react-router-dom";
import PrimaryButton from "../components/buttons/PrimaryButton";

export default function LandingPage() {
    const navigate = useNavigate();

    return (
        <main
            className="
                min-h-[100dvh]
                w-full
                bg-surface-base
                bg-cover
                bg-center
                bg-no-repeat
                flex
                items-center
                justify-center
                px-5
                py-8
            "
            style={{
                backgroundImage: "url('/images/landing-bg.png')",
            }}
        >
            <div
                className="
                    flex
                    w-full
                    max-w-sm
                    flex-col
                    items-center
                    text-center
                    sm:max-w-md
                    lg:max-w-xl
                "
            >
                <div className="flex flex-col items-center">
                    <img
                        src="/logos/atlas-logo-landing.png"
                        alt="ATLAS"
                        className="mb-4 w-52 sm:w-60 md:w-72 lg:w-80"
                    />
                </div>
                <div 
                    className="
                        flex
                        items-center
                        justify-center
                        gap-4
                        text-heading-small
                        uppercase
                        text-white
                        sm:gap-6
                        sm:text-base
                        md:gap-8
                        md:text-lg
                        py-4
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
                        max-w-xs
                        flex-col
                        gap-2
                        sm:max-w-sm
                        md:gap-5
                        lg:max-w-md
                        items-center
                    "
                >
                    <PrimaryButton size="medium" onClick={() => navigate("/login")}>
                      Login
                    </PrimaryButton>

                    <PrimaryButton size="medium" onClick={() => navigate("/game")}>
                      Guest Mode
                    </PrimaryButton>
                </div>

                <button
                    type="button"
                    className="
                        mt-7
                        md:mt-7
                        cursor-pointer
                        text-body-primary
                        uppercase
                        tracking-[0.08em]
                        text-white
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
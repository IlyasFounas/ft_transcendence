import { useNavigate } from "react-router-dom";
import PrimaryButton from "../components/buttons/PrimaryButton";

export default function LandingPage() {
    const navigate = useNavigate();

    return (
        <main
            className="
                min-h-screen
                w-full
                bg-surface-base
                bg-cover
                bg-center
                bg-no-repeat
                flex
                items-center
                justify-center
                px-6
            "
            style={{
                backgroundImage: "url('/images/landing-bg.png')",
            }}
        >
            <div
                className="
                    flex
                    w-full
                    max-w-xl
                    flex-col
                    items-center
                    text-center
                "
            >
                <div className="mb-8 flex flex-col items-center">
                    <img
                        src="/logos/atlas-logo-square.png"
                        alt="ATLAS"
                        className="mb-4 w-64 md:w-80 lg:w-120"
                    />
                </div>

                <div className="flex w-full max-w-sm flex-col gap-4">
                    <PrimaryButton size="large" onClick={() => navigate("/login")}>
                      Login
                    </PrimaryButton>

                    <PrimaryButton size="medium" onClick={() => navigate("/game")}>
                      Start Guest Mode
                    </PrimaryButton>
                </div>

                <button
                    type="button"
                    className="
                        mt-7
                        cursor-pointer
                        text-sm
                        font-medium
                        tracking-[0.08em]
                        text-white
                        underline
                        underline-offset-4
                        transition
                        hover:text-primary
                    "
                >
                    How to play
                </button>
            </div>
        </main>
    );
}
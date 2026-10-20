import AuthenticationCard from "../components/cards/AuthenticationCard";
import AuthenticationLogin from "../components/authentication/AuthenticationLogin";
import { useState } from "react";
import AuthenticationWelcome from "../components/authentication/AuthenticationWelcome";
import AuthenticationSignUp from "../components/authentication/AuthenticationSignUp";
import { useNavigate } from "react-router-dom";

type AuthMode = "options" | "signup" | "login";

export default function AuthenticationPage () {
    
    const [mode, setMode] = useState<AuthMode>("options");

    const navigate = useNavigate();

    const handleGoogleLogin = () => {
        window.location.href = "/api/auth/google";
    };
    
    return (
        <main
            className="
                flex
                min-h-screen
                w-full
                flex-col
                items-center
                justify-center
                gap-6

                bg-surface-base
                bg-cover
                bg-center
                bg-no-repeat
                bg-[url('/backgrounds/landing_background.png')]
                px-8
                md:px-10
                md:py-5
                lg:flex-row
            "
        >
            <div 
                className="
                    hidden
                    w-full 
                    max-w-lg 
                    items-center
                    justify-center
                    sm:flex
                "
            >
                <button
                    type="button"
                    onClick={() => navigate("/")}
                    className="cursor-pointer"
                    aria-label="Go to landing page"
                >
                    <img
                        src="/logos/atlas-logo-landing.png"
                        alt="ATLAS"
                        className="w-72 md:w-80 lg:w-116"
                    />
                </button>
            </div>
            <AuthenticationCard>
                {mode === "options" && (
                    <AuthenticationWelcome
                        onGoogleLogin={handleGoogleLogin}
                        onEmailSignUp={() => setMode("signup")}
                        onEmailLogin={() => setMode("login")}
                    />
                )}

                {mode === "signup" && (
                    <AuthenticationSignUp
                        onLogin={() => setMode("login")}
                    />
                )}

                {mode === "login" && (
                    <AuthenticationLogin
                        onBackToOptions={() => setMode("options")}
                    />
                )}
            </AuthenticationCard>
        </main>
    );
}
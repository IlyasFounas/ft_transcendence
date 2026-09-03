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

    return (
        <main
            className="
                flex
                min-h-screen
                w-full
                items-center
                justify-center
                gap-20
                bg-surface-base
                bg-cover
                bg-center
                bg-no-repeat
                px-6
            "
            style={{
                backgroundImage: "url('/images/landing-bg.png')",
            }}
        >
            <div className="flex w-full max-w-lg items-center justify-center">
                <button
                    type="button"
                    onClick={() => navigate("/")}
                    className="cursor-pointer"
                    aria-label="Go to landing page"
                >
                    <img
                        src="/logos/atlas-logo-square.png"
                        alt="ATLAS"
                        className="w-72 md:w-80 lg:w-96"
                    />
                </button>
            </div>
            <AuthenticationCard>
                {mode === "options" && (
                    <AuthenticationWelcome
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
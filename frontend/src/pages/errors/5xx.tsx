import { useNavigate } from "react-router-dom";
import PrimaryButton from "../../components/buttons/PrimaryButton";

export default function ServerErrorPage() {
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
            <div className="flex w-full max-w-2xl flex-col items-center text-center">
                <img
                    src="/logos/atlas-logo-landing.png"
                    alt="ATLAS"
                    className="
                        w-32
                        md:w-48
                        lg:w-72
                        mb-8
                    "
                />

                <img
                    src="/images/5XX.png"
                    alt="server error"
                    className="
                        w-96
                        md:w-120
                        lg:w-240
                        h-auto
                    "
                />
                <p className="mt-4  max-w-md text-body-primary md:text-heading-medium text-text-secondary uppercase">
                   Something went wrong. Please try again later.
                </p>
                <div className="
                    flex
                    w-full
                    flex-col
                    items-center
                    mt-6
                    max-w-xs
                    lg:hidden
                    "
                >
                    <PrimaryButton
                        size="small"
                        onClick={() => navigate("/")}
                    >
                        Return to Home
                    </PrimaryButton>
                </div>
                <div className="
                    hidden
                    lg:block
                    w-full
                    flex-col
                    items-center
                    mt-6
                    max-w-xs
                    lg:max-w-md
                    "
                >
                    <PrimaryButton
                        size="medium"
                        onClick={() => navigate("/")}
                    >
                        Return to Home
                    </PrimaryButton>
                </div>
            </div>
        </main>
    );
}
import { useNavigate } from "react-router-dom";
import PrimaryButton from "../../components/buttons/PrimaryButton";

export default function NotFoundPage() {
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
            <div className="flex flex-col">
                <p className="text-[180px] font-semibold">
                    404
                </p>
                <p className="text-heading-large">
                    Page Not Found
                </p>
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
                    onClick={() => navigate("/")}
                >
                    Go to landing page
                </PrimaryButton>
            </div>

        </main>
    );
}
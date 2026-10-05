import { useState } from "react";
import { tutorialPages } from "../../data/tutorial";
import PrimaryButton from "../buttons/PrimaryButton";
import { X } from "lucide-react";

interface TutorialCarouselProps {
    mode: "landing" | "game";
    onClose: () => void;
    onStartOperation?: () => void;
}

export default function TutorialCarousel({
    mode,
    onClose,
    onStartOperation,
}: TutorialCarouselProps) {
    const [currentPage, setCurrentPage] = useState(0);

    const page = tutorialPages[currentPage];
    const isLastPage = currentPage === tutorialPages.length - 1;

    const handleNext = () => {
        if (!isLastPage) {
            setCurrentPage((current) => current + 1);
        }
    };

    const handleBack = () => {
        if (currentPage > 0) {
            setCurrentPage((current) => current - 1);
        }
    };

    const handleLastPageAction = () => {
        if (mode === "landing") {
            onStartOperation?.();
        } else {
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4">
            <div
                className="
                    relative
                    w-full
                    max-w-208
                    h-157
                    rounded-xl
                    border
                    border-border-default
                    bg-surface-base
                    p-6
                    md:p-8
                    flex
                    flex-col
                "
            >
                <div className="flex items-center justify-between">
                    <img
                        src="/logos/atlas-logo.png"
                        alt="Atlas"
                        className="h-12"
                    />
                    <X 
                        size={24}
                        type="button"
                        onClick={onClose}
                        className="text-text-primary transition hover:text-primary"
                    />
                </div>

                <div className="mt-4 flex items-start justify-between gap-4">
                    <div className="text-left">
                        <h2 className="text-hero uppercase">
                            <span className="text-text-primary">
                                {page.title}{" "}
                            </span>
                            <span className="text-primary">
                                {page.highlight}
                            </span>
                        </h2>
                        <div className="h-12">
                            <p className="mt-2 text-heading-small text-text-primary">
                                {page.description}
                            </p>
                        </div>
                    </div>

                    <span className="shrink-0 text-text-secondary text-heading-small">
                        {String(currentPage + 1).padStart(2, "0")} /{" "}
                        {String(tutorialPages.length).padStart(2, "0")}
                    </span>
                </div>

                <div className="relative mt-3 h-87.5 w-full shrink-0">
                    <div className="flex h-full w-full items-center justify-center">
                        <img
                            src={page.image}
                            alt={page.title}
                            className="max-h-full max-w-full object-contain"
                        />
                    </div>

                    {isLastPage && (
                       <div className="absolute bottom-12 left-1/2 w-60 -translate-x-1/2 translate-y-1/2">
                            <PrimaryButton
                                type="button"
                                size="small"
                                onClick={handleLastPageAction}
                            >
                                {mode === "landing"
                                    ? "START OPERATION"
                                    : "CONTINUE TO PLAY"}
                            </PrimaryButton>
                        </div>
                    )}
                </div>

                <div className="mt-auto flex items-center justify-between">
                    {currentPage > 0 ? (
                        <button
                            type="button"
                            onClick={handleBack}
                            className="
                                text-button
                                uppercase
                                text-heading-small
                                text-text-secondary
                                transition
                                hover:text-primary
                            "
                        >
                            ← BACK
                        </button>
                    ) : (
                        <div className="w-17.5" />
                    )}

                    <div className="flex items-center gap-2">
                        {tutorialPages.map((_, index) => (
                            <span
                                key={index}
                                className={`
                                    h-2
                                    w-2
                                    rounded-full
                                    ${
                                        index === currentPage
                                            ? "bg-primary"
                                            : "bg-text-secondary"
                                    }
                                `}
                            />
                        ))}
                    </div>

                    {!isLastPage ? (
                        <button
                            type="button"
                            onClick={handleNext}
                            className="
                                text-button
                                uppercase
                                text-heading-small
                                text-primary
                                transition
                                hover:text-text-primary
                            "
                        >
                            NEXT →
                        </button>
                    ) : (
                        <div className="w-17.5" />
                    )}
                </div>
            </div>
        </div>
    );
}
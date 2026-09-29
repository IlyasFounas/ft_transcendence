import { useState } from "react";

interface TutorialCarouseProps {
    images: string[];
}

export default function TutorialCarousel({
    images,
}: TutorialCarouseProps) {

    const [currentIndex, setCurrentIndex] = useState(0);

    const nextSlide = () => {
        setCurrentIndex((prevIndex) => (prevIndex + 1) % images.length);
    };

    const prevSlide = () => {
        setCurrentIndex((prevIndex) => (prevIndex - 1 + images.length) % images.length);
    };

    return (
        <div className="relative w-full max-w-3xl mx-auto">
            <div className="overflow-hidden relative h-150">
                {images.map((image, index) => (
                    <div
                        key={index}
                        className={`absolute inset-0 transition-transform 
                            ${index === currentIndex ? 'translate-x-0' : 'translate-x-full'}`}
                    >
                        <img 
                            src={image}
                            alt={`Tutorial slide ${index + 1}`}
                            className="w-full h-full object-cover"
                        />
                    </div>
                ))}
            </div>
            <button 
                className="absolute top-1/2 left-0 transform -translate-y-1/2 bg-gray-800 text-white p-2"
                onClick={prevSlide}
            >
                Prev
            </button>
            <button 
                className="absolute top-1/2 right-0 transform -translate-y-1/2 bg-gray-800 text-white p-2"
                onClick={nextSlide}
            >
                Next
            </button>
        </div>
    );
}
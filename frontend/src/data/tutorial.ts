export interface TutorialPage {
    title: string;
    highlight: string;
    description: string;
    image: string;
    moibileImage: string;
}

export const tutorialPages: TutorialPage[] = [
    {
        title: "INVESTIGATE",
        highlight: "THE SIGNAL",
        description: "Each operation begins with an unknown visual transmission. Explore the location and determine where the signal was captured before the signal is lost.",
        image: "/images/tutorial/tutorial-1.png",
        moibileImage: "/images/tutorial/tutorial-1 (phone).png",
    },
    {
        title: "SEARCH",
        highlight: "FOR CLUE",
        description: "Every detail can reveal the location.",
        image: "/images/tutorial/tutorial-2.png",
        moibileImage: "/images/tutorial/tutorial-2 (phone).png",
    },
    {
        title: "SUBMIT",
        highlight: "YOUR GUESS",
        description: "When you are ready, switch to the map and place a pin where you think the signal was captured.",
        image: "/images/tutorial/tutorial-3.png",
        moibileImage: "/images/tutorial/tutorial-3 (phone).png",
    },
    {
        title: "SCORE",
        highlight: "MORE XP",
        description: "Every detail can reveal the location.",
        image: "/images/tutorial/tutorial-4.png",
        moibileImage: "/images/tutorial/tutorial-4 (phone).png",
    },
    {
        title: "READY TO",
        highlight: "OPERATE?",
        description: "Use your observation skills, explore the world, and submit your best guess before the signal is lost.",
        image: "/images/tutorial/tutorial-4.png",
        moibileImage: "/images/tutorial/tutorial-5 (phone).png",
    },
];
type Props = {
    countdown: number;
};

export default function Countdown({ countdown }: Props)
{
	function formatTime(value: number)
	{
		const minutes = Math.floor(value / 60);
		const seconds = value % 60;

		return `${minutes}:${seconds.toString().padStart(2, "0")}`;
	}

	return (
		<div className={`
			w-[75px] h-[40px] shrink-0
			lg:w-[200px] lg:h-[75px]
			flex items-center justify-center
			rounded-lg
			px-6 py-4
			bg-surface-base
			outline outline-1
			shadow-[inset_0px_0px_20px_0px_rgba(0,175,255,0.05)]
			uppercase
			text-heading-small lg:text-hero 
			${countdown <= 10 ? "outline-red-500" : "outline-sky-950"}
		`}>
			<div className={` transition-colors duration-300 ${countdown <= 10 ? "text-red-500 animate-pulse" : "text-white"} `}>
				0{formatTime(countdown)}
			</div>
		</div>
	)
}
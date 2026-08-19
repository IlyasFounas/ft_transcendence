export default function SignalIntegrityCard({ countDown, totalTime }: {countDown : number, totalTime : number}) 
{
	const TOTAL_SIGNALS = 15;

	const percentage = Math.round((countDown / totalTime) * 100);
	const activeSignals = Math.ceil((percentage / 100) * TOTAL_SIGNALS);

	return (
		<div className="
			w-full h-[75px]
			flex flex-col justify-center
			rounded-lg
			px-2
			shadow-[inset_0px_0px_20px_0px_rgba(0,175,255,0.05)] outline outline-1 outline-sky-950
			uppercase
 			">
			<h1 className="
				text-heading-small">Signal Integrity</h1>
			<div className="flex flex-row justify-between items-center">
				<div className="flex gap-1">
					{Array.from({ length: TOTAL_SIGNALS }).map((_, index) => (
						<div
							key={index}
							className={`
								w-3 h-6 rounded-xs
								${index < activeSignals
									? "bg-cyan-400"
									: "bg-sky-950"}
							`}
						/>
					))}
				</div>
				<div className="text-heading-large text-primary">{percentage}%</div>

			</div>
		</div>
	)
}
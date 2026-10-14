interface MissionSettingsCardProps {
  rounds: number;
  timeLimit: number;
};

export default function MissionSettingsCard({rounds, timeLimit}: MissionSettingsCardProps)
{
	function formatTime(value: number)
	{
		if (value < 60)
			return `${value} sec`;

		const minutes = Math.floor(value / 60);
		const seconds = value % 60;

		if (seconds === 0)
			return `${minutes} min`;

		return `${minutes} min ${seconds}`;
	}
	return (
		<div 
		className="rounded-[12px] bg-surface-base
			border border-border-default/30
			shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden">
	  		<div className="
				flex 
				items-center 
				justify-between 
				border-b 
				border-border-default/30
				bg-gradient-to-t from-black to-bg-surface-base 
				bg-surface-title 
				px-4 py-3">
			<span className="text-heading-small text-text-primary/70 uppercase ">
					Mission Settings
			</span>
	  		</div>
			<div className="w-full min-h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center justify-between px-[16px] 
			text-heading-small uppercase">
				<div>Rounds</div>
				<div>{rounds}</div>
			</div>
			<div className="w-full min-h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center justify-between px-[16px] 
			text-heading-small uppercase">
				<div>Time Limit</div>
				<div>{formatTime(timeLimit)}</div>
			</div>
   		</div>
	)
}
import { CircleMinus } from 'lucide-react';
import { CirclePlus } from 'lucide-react';


interface MissionSettingsCardProps {
	className?: string;
};

export default function MissionSettingsCard({className = ""}: MissionSettingsCardProps)
{
	return (
		<div 
		className={`rounded-[12px] bg-surface-base
			border border-border-default/30
			shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden ${className}`}>
	  		<div className="flex items-center justify-between border-b border-border-default/30 
				  bg-surface-title px-4 py-3">
			<span className="text-heading-small text-text-primary/70 uppercase">
					Mission Settings
			</span>
	  		</div>
			<div className="w-full min-h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center justify-between px-[16px] 
			text-heading-small uppercase">
				<div>Rounds</div>
				<div className="flex flex-row justify-between gap-4 shrink-0 w-[120px]">
					<CircleMinus></CircleMinus>
					<div>3</div>
					<CirclePlus></CirclePlus>
				</div>
			</div>
			<div className="w-full min-h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center justify-between px-[16px] 
			text-heading-small uppercase">
				<div>Time Limit</div>
				<div className="flex flex-row gap-4">
					<CircleMinus></CircleMinus>
					<div>2 min</div>
					<CirclePlus></CirclePlus>
				</div>
			</div>
   		</div>
	)
}
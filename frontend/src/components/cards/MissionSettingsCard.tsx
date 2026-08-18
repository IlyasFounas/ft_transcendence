import {useState} from 'react';
import { Minus } from 'lucide-react';
import { Plus } from 'lucide-react';


interface MissionSettingsCardProps {
	className?: string;
};



export default function MissionSettingsCard({className = ""}: MissionSettingsCardProps)
{
	const [rounds, setRounds] = useState(5);
	const [timeLimit, setTimeLimit] = useState(2.5);

	function formatTime(value : number)
	{
		if (value % 1 === 0.5)
			 return `${Math.floor(value)} min 30`;

   		return `${value} min`;
	}


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
				<div className="flex flex-row justify-between gap-4 shrink-0 w-[150px]">
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (rounds > 1) {
								setRounds(rounds - 1);
							}}}>
						<Minus/>
					</button>
					<div>{rounds}</div>
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (rounds < 10) {
								setRounds(rounds + 1);
							}}}>
						<Plus/>
					</button>
				</div>
			</div>
			<div className="w-full min-h-[65px]
			outline outline-1 outline-offset-[-1px] outline-cyan-800/30
			flex items-center justify-between px-[16px] 
			text-heading-small uppercase">
				<div>Time Limit</div>
				<div className="flex flex-row justify-between gap-4 shrink-0 w-[150px]">
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (timeLimit > 1) {
								setTimeLimit(timeLimit - 0.5);
							}}}>
						<Minus/>
					</button >
					<div>{formatTime(timeLimit)}</div>
					<button className="
						cursor-pointer
						transition-transform
						duration-200
						ease-out
						hover:scale-[1.2]
						active:scale-[0.98]
						"
						onClick={() => {
							if (timeLimit < 5) {
								setTimeLimit(timeLimit + 0.5);
							}}}>
						<Plus/>
					</button>
				</div>
			</div>
   		</div>
	)
}
import { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface RoomCodeCardProps {
	code: string;
};

export default function RoomCodeCard({ code}: RoomCodeCardProps)
{
	const [copied, setCopied] = useState(false);

	const  handleCopy = async () => {
		await navigator.clipboard.writeText(code);
		setCopied(true);

		setTimeout(() => {
        	setCopied(false);
    	}, 2000);
	}

	return (
		<div 
		className="
			rounded-[12px] 
			bg-surface-base
			border border-border-default/30
			shadow-[0_0_10px_rgba(40,215,255,0.5)] 
			overflow-hidden">
      		<div className="
				flex 
				items-center 
				justify-between 
				border border-border-default/30
				bg-gradient-to-t from-black to-bg-surface-base
				px-4 py-3">
        	<span className="text-heading-small text-text-primary/70 uppercase">
        		Room Code
        	</span>
      		</div>
			<div className="p-4 flex justify-between text-heading-large items-center">
				<span>
					{code}
				</span>
				<div className="flex items-center gap-2">
					<button
						type="button"
						onClick={handleCopy}
						className="
							cursor-pointer
							transition-transform
							duration-200
							ease-out
							hover:scale-[1.2]
							active:scale-[0.98]"
						>
						{copied ? <Check/> : <Copy/>}
					</button>
				</div>
			</div>
   		</div>
	)
}
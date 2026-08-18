import { useState } from 'react';
import { Copy } from 'lucide-react';

interface RoomCodeCardProps {
	code: string;
	className?: string;
};

export default function RoomCodeCard({ code = "ATL-7961", className = ""}: RoomCodeCardProps)
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
		className={`rounded-[12px] bg-surface-base
			border border-border-default/30
            shadow-[0_0_10px_rgba(40,215,255,0.5)] overflow-hidden ${className}`}>
      		<div className="flex items-center justify-between border-b border-border-default/30 
                  bg-surface-title px-4 py-3">
        	<span className="text-heading-small text-text-primary/70 uppercase">
        		Room Code
        	</span>
      		</div>
			<div className="p-4 flex justify-between text-heading-large items-center">
				<span>
					{code}
				</span>
				<div className="flex items-center gap-2">
					{copied && (
						<span className="text-sm text-text-primary">
							Copied!
						</span>
					)}

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
						<Copy />
					</button>
				</div>
			</div>
   		</div>
	)
}
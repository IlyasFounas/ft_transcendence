import { Copy } from 'lucide-react';

interface RoomCodeCardProps {
	code: string;
	className?: string;
};

export default function RoomCodeCard({ code = "ATL-7961", className = ""}: RoomCodeCardProps)
{
	return (
		<div 
		className={`rounded-[12px]
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
				<Copy></Copy>
			</div>
   		</div>
	)
}
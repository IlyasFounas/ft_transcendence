interface ChatBoxCardProps {
	variant: "received" | "sent";
	message: string;
	time: string;
}

const variants = {
	sent: {
		wrapper: "justify-end",
		bubble: "bg-surface-elevated items-end",
	},
	received: {
		wrapper: "justify-start",
		bubble: "bg-neutral-700 items-start"
	},
};

export default function ChatBoxCard({
	variant,
	message,
	time,
}: ChatBoxCardProps){
	
	const style = variants[variant];

	return(
		<div className={`w-full px-[5px] flex
    		${style.wrapper}`}
		>
			<div className={`max-w-64 px-2.5 py-[5px] rounded-[10px] inline-flex flex-col ${style.bubble}`}>
				<div className="min-w-0 max-w-full wrap-break-words">
					<div className="text-body-primary wrap-break-words">{message}</div>
				</div>
				<div className="self-stretch flex justify-end">
					<div className="text-label-small text-text-secondary ">{time}</div>
				</div>
			</div>
		</div>
	)
}
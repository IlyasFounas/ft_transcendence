interface StatusBadgeProps {
  status: "ready" | "waiting";
}

const variants = {
	ready: {
		text: "text-status-green",
		dot: "bg-status-green",
	},
	waiting: {
		text: "text-status-yellow",
		dot: "bg-status-yellow"
	},
};

export default function StatusBadge({ status }: StatusBadgeProps) {

	const style = variants[status];

  	return (
		<div className="inline-flex items-center gap-2.5">
  			<span className={`text-heading-small uppercase ${style.text}`}>{status}</span>
  			<div className={`size-1.5 rounded-full ${style.dot}`}></div>
		</div>
  );
}
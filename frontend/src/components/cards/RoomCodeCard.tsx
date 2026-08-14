import { Copy } from 'lucide-react';

export default function RoomCodeCard({ code = "ATL-7961"}: { code: string;})
{
	return (
		<div>
			<h1>Room Code</h1>
			<div>
				{code}
				<Copy></Copy>
			</div>
		</div>
	)
}

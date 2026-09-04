import { useState, useEffect } from "react";

export default function WaitingDots() {
  const [dotCount, setDotCount] = useState(1);

  useEffect(() => {
	const interval = setInterval(() => {
	  setDotCount((prev) => (prev >= 3 ? 0 : prev + 1));
	}, 500);

	return () => clearInterval(interval);
  }, []);

  return (
	<span>
	  Waiting for others
	  <span className="inline-block w-4 text-left">{".".repeat(dotCount)}</span>
	</span>
  );
}
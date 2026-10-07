import { useState, useEffect, useRef } from "react";

export function useCountUp(targetValue: number, duration = 800) {
  const [displayValue, setDisplayValue] = useState(0);
  const [justFinished, setJustFinished] = useState(false);
  const previousValue = useRef(0);

  useEffect(() => {
    const startValue = previousValue.current;
    const diff = targetValue - startValue;

    if (diff === 0) return;

    const startTime = performance.now();

    const animate = (now: number) => {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);

      setDisplayValue(Math.round(startValue + diff * eased));

      if (progress < 1) {
        requestAnimationFrame(animate);
      } else {
        previousValue.current = targetValue;
        setJustFinished(true);
        setTimeout(() => setJustFinished(false), 300);
      }
    };

    requestAnimationFrame(animate);
  }, [targetValue, duration]);

  return { displayValue, justFinished };
}
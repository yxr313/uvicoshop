"use client";

import { useEffect, useState } from "react";
import { animate } from "motion/react";
import { formatUsd } from "@uvico/pricing";

interface AnimatedPriceProps {
  value: number;
}

export function AnimatedPrice({ value }: AnimatedPriceProps) {
  const [displayValue, setDisplayValue] = useState(value);

  useEffect(() => {
    const controls = animate(displayValue, value, {
      duration: 0.5,
      ease: [0.22, 1, 0.36, 1], // Nuestro easing Uvico
      onUpdate: (latest) => setDisplayValue(latest),
    });
    return () => controls.stop();
  }, [value]);

  return (
    <span className="text-[2.5rem] leading-none font-display font-medium tracking-tighter text-foreground">
      {formatUsd(displayValue)}
    </span>
  );
}

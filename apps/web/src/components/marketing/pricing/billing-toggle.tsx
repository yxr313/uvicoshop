"use client";

import { cn } from "@/lib/utils";

interface BillingToggleProps {
  cycle: "monthly" | "annual";
  onChange: (cycle: "monthly" | "annual") => void;
}

export function BillingToggle({ cycle, onChange }: BillingToggleProps) {
  return (
    <div className="flex items-center justify-center gap-6 mb-16">
      <button
        onClick={() => onChange("monthly")}
        className={cn(
          "text-body transition-colors pb-1 border-b-2",
          cycle === "monthly" 
            ? "text-foreground font-medium border-foreground" 
            : "text-muted-foreground border-transparent hover:text-foreground"
        )}
      >
        Mensual
      </button>
      <button
        onClick={() => onChange("annual")}
        className={cn(
          "text-body transition-colors pb-1 border-b-2 flex items-center gap-2",
          cycle === "annual" 
            ? "text-foreground font-medium border-foreground" 
            : "text-muted-foreground border-transparent hover:text-foreground"
        )}
      >
        Anual
        <span className="bg-accent text-accent-foreground text-[10px] font-bold px-2 py-0.5 rounded-full">
          Ahorra 17%
        </span>
      </button>
    </div>
  );
}

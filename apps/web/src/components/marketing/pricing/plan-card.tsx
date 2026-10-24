"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Check, Zap } from "lucide-react";
import { formatUsd } from "@uvico/pricing";
import { cn } from "@/lib/utils";

interface PlanCardProps {
  plan: {
    id: string;
    name: string;
    tagline: string;
    isFounder: boolean;
    features: string[];
  };
  price: {
    effectiveMonthlyUsd: number;
    amountUsd: number;
  };
  cycle: "monthly" | "annual";
  isPopular?: boolean;
}

export function PlanCard({ plan, price, cycle, isPopular }: PlanCardProps) {
  return (
    <div 
      className={cn(
        "relative flex flex-col p-8 rounded-[2rem] border transition-all duration-500 hover:shadow-xl",
        plan.isFounder ? "border-foreground shadow-lg bg-white" : "border-border bg-neutral-50",
        isPopular && "border-accent ring-1 ring-accent"
      )}
    >
      {plan.isFounder && (
        <div className="absolute top-0 inset-x-0 h-1.5 bg-accent rounded-t-[2rem]" />
      )}
      
      {isPopular && (
        <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-accent text-accent-foreground text-xs font-bold px-4 py-1 rounded-full flex items-center gap-1 shadow-sm">
          <Zap size={12} className="fill-current" />
          Más popular
        </div>
      )}

      <div className="mb-8">
        <h3 className="text-heading-lg font-display text-foreground mb-2">{plan.name}</h3>
        <p className="text-body-sm text-muted-foreground h-10">{plan.tagline}</p>
      </div>

      <div className="mb-8 h-20">
        <div className="flex items-baseline gap-1">
          <span className="text-[3rem] leading-none font-display font-medium tracking-tighter text-foreground transition-all duration-300">
            {formatUsd(price.effectiveMonthlyUsd)}
          </span>
          <span className="text-body text-muted-foreground">/mes</span>
        </div>
        {cycle === "annual" && (
          <p className="text-sm text-muted-foreground mt-2">
            Facturado {formatUsd(price.amountUsd)} al año
          </p>
        )}
      </div>

      <div className="flex-1 space-y-4 mb-8">
        {plan.features.map((feature, idx) => (
          <div key={idx} className="flex items-start gap-3">
            <Check size={18} className="text-foreground mt-0.5 shrink-0" />
            <span className="text-body-sm text-foreground">{feature}</span>
          </div>
        ))}
      </div>

      <Link href={`/registro?plan=${plan.id}`} className="w-full mt-auto">
        <Button variant={plan.isFounder ? "default" : "dark"} className="w-full">
          Empezar ahora
        </Button>
      </Link>
    </div>
  );
}

export type PlanId = "fundador" | "esencial" | "profesional" | "empresa";
export type BillingCycle = "monthly" | "annual";

export const PLAN_IDS: readonly PlanId[] = [
  "fundador",
  "esencial",
  "profesional",
  "empresa",
] as const;

export const BILLING_CYCLES: readonly BillingCycle[] = ["monthly", "annual"] as const;

export const ANNUAL_DISCOUNT_RATE = 0.17;
export const MONTHS_PER_YEAR = 12;

export interface PlanDefinition {
  readonly id: PlanId;
  readonly name: string;
  readonly tagline: string;
  readonly monthlyUsd: number;
  readonly isFounder: boolean;
}

export const PLANS = {
  fundador: {
    id: "fundador",
    name: "Fundador",
    tagline: "Para emprendedores que arrancan hoy",
    monthlyUsd: 9,
    isFounder: true,
  },
  esencial: {
    id: "esencial",
    name: "Esencial",
    tagline: "Para tiendas que ya venden",
    monthlyUsd: 19.99,
    isFounder: false,
  },
  profesional: {
    id: "profesional",
    name: "Profesional",
    tagline: "Para comercios que escalan",
    monthlyUsd: 49.99,
    isFounder: false,
  },
  empresa: {
    id: "empresa",
    name: "Empresa",
    tagline: "Para operaciones con equipo y multi-tienda",
    monthlyUsd: 99.99,
    isFounder: false,
  },
} as const satisfies Record<PlanId, PlanDefinition>;

export interface PlanPriceResult {
  readonly planId: PlanId;
  readonly cycle: BillingCycle;
  readonly amountUsd: number;
  readonly effectiveMonthlyUsd: number;
  readonly savingsVsMonthlyUsd: number;
  readonly savingsRate: number;
}

export function roundUsd(amount: number): number {
  return Math.round(amount * 100) / 100;
}

export function getPlanPrice(planId: PlanId, cycle: BillingCycle): PlanPriceResult {
  const plan = PLANS[planId];
  const monthlyTotal = roundUsd(plan.monthlyUsd * MONTHS_PER_YEAR);

  if (cycle === "monthly") {
    return {
      planId,
      cycle,
      amountUsd: plan.monthlyUsd,
      effectiveMonthlyUsd: plan.monthlyUsd,
      savingsVsMonthlyUsd: 0,
      savingsRate: 0,
    };
  }

  const amountUsd = roundUsd(monthlyTotal * (1 - ANNUAL_DISCOUNT_RATE));

  return {
    planId,
    cycle,
    amountUsd,
    effectiveMonthlyUsd: roundUsd(amountUsd / MONTHS_PER_YEAR),
    savingsVsMonthlyUsd: roundUsd(monthlyTotal - amountUsd),
    savingsRate: ANNUAL_DISCOUNT_RATE,
  };
}

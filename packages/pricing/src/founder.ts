import { PLANS, roundUsd, type PlanId } from "./plans";

export const FOUNDER_PLAN_ID: PlanId = "fundador";

export const FOUNDER_TRIAL_DAYS = 7;

export const FOUNDER_INTRO_PRICE_USD = 1;
export const FOUNDER_INTRO_MONTHS = 3;

export const FOUNDER_SLOTS_TOTAL = 500;

export interface FounderIntroSummary {
  readonly introPriceUsd: number;
  readonly introMonths: number;
  readonly introTotalUsd: number;
  readonly regularMonthlyUsd: number;
  readonly trialDays: number;
}

export function getFounderIntroSummary(): FounderIntroSummary {
  return {
    introPriceUsd: FOUNDER_INTRO_PRICE_USD,
    introMonths: FOUNDER_INTRO_MONTHS,
    introTotalUsd: roundUsd(FOUNDER_INTRO_PRICE_USD * FOUNDER_INTRO_MONTHS),
    regularMonthlyUsd: PLANS.fundador.monthlyUsd,
    trialDays: FOUNDER_TRIAL_DAYS,
  };
}

export function getFounderSlotsRemaining(slotsTaken: number): number {
  if (!Number.isInteger(slotsTaken) || slotsTaken < 0) {
    throw new Error("slotsTaken must be a non-negative integer");
  }
  return Math.max(0, FOUNDER_SLOTS_TOTAL - slotsTaken);
}

export function isFounderSlotAvailable(slotsTaken: number): boolean {
  return getFounderSlotsRemaining(slotsTaken) > 0;
}

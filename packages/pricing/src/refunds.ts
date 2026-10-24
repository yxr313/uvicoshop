import { roundUsd } from "./plans";

export type RefundReason = "changed-mind" | "merchant-fault" | "chargeback";

export interface StoreCreditOutcome {
  readonly type: "store-credit";
  readonly amountUsd: number;
  readonly absorbedBy: "merchant";
  readonly useScope: "any-uvico-store";
}

export interface CashRefundOutcome {
  readonly type: "cash-refund";
  readonly amountUsd: number;
  readonly absorbedBy: "merchant";
  readonly trigger: "merchant-fault";
}

export interface ChargebackDebitOutcome {
  readonly type: "chargeback-debit";
  readonly amountUsd: number;
  readonly absorbedBy: "merchant";
  readonly recoverableFrom: "retained-balance" | "future-payouts";
}

export type RefundOutcome = StoreCreditOutcome | CashRefundOutcome | ChargebackDebitOutcome;

export function resolveRefund(reason: RefundReason, amountUsd: number): RefundOutcome {
  if (!Number.isFinite(amountUsd) || amountUsd <= 0) {
    throw new Error("amountUsd must be a positive finite number");
  }

  const rounded = roundUsd(amountUsd);

  switch (reason) {
    case "changed-mind":
      return {
        type: "store-credit",
        amountUsd: rounded,
        absorbedBy: "merchant",
        useScope: "any-uvico-store",
      };
    case "merchant-fault":
      return {
        type: "cash-refund",
        amountUsd: rounded,
        absorbedBy: "merchant",
        trigger: "merchant-fault",
      };
    case "chargeback":
      return {
        type: "chargeback-debit",
        amountUsd: rounded,
        absorbedBy: "merchant",
        recoverableFrom: "retained-balance",
      };
  }
}

import { roundUsd } from "./plans";

export const DROPSHIPPING_COMMISSION_RATE = 0.05;

export type PayoutFrequency = "weekly" | "biweekly";

export const PAYOUT_FREQUENCIES: readonly PayoutFrequency[] = ["weekly", "biweekly"] as const;

export const DEFAULT_PAYOUT_FREQUENCY: PayoutFrequency = "weekly";

export const GATEWAY_LIQUIDATION_MIN_DAYS = 1;
export const GATEWAY_LIQUIDATION_MAX_DAYS = 3;

export interface SaleSplitInput {
  readonly saleAmountUsd: number;
  readonly supplierCostUsd: number;
  readonly gatewayFeeUsd: number;
}

export interface SaleSplit {
  readonly grossUsd: number;
  readonly supplierUsd: number;
  readonly uvicoCommissionUsd: number;
  readonly gatewayFeeUsd: number;
  readonly merchantNetUsd: number;
}

export function calculateSaleSplit(input: SaleSplitInput): SaleSplit {
  const { saleAmountUsd, supplierCostUsd, gatewayFeeUsd } = input;

  if (!Number.isFinite(saleAmountUsd) || saleAmountUsd <= 0) {
    throw new Error("saleAmountUsd must be a positive finite number");
  }
  if (!Number.isFinite(supplierCostUsd) || supplierCostUsd < 0) {
    throw new Error("supplierCostUsd must be a non-negative finite number");
  }
  if (!Number.isFinite(gatewayFeeUsd) || gatewayFeeUsd < 0) {
    throw new Error("gatewayFeeUsd must be a non-negative finite number");
  }
  if (supplierCostUsd + gatewayFeeUsd > saleAmountUsd) {
    throw new Error("supplierCostUsd + gatewayFeeUsd exceed saleAmountUsd");
  }

  const uvicoCommissionUsd = roundUsd(saleAmountUsd * DROPSHIPPING_COMMISSION_RATE);
  const merchantNetUsd = roundUsd(
    saleAmountUsd - supplierCostUsd - uvicoCommissionUsd - gatewayFeeUsd,
  );

  return {
    grossUsd: saleAmountUsd,
    supplierUsd: supplierCostUsd,
    uvicoCommissionUsd,
    gatewayFeeUsd,
    merchantNetUsd,
  };
}

export interface SaleMarginValidation {
  readonly ok: boolean;
  readonly reason?: string;
}

export function validateSaleMargin(input: SaleSplitInput): SaleMarginValidation {
  try {
    const split = calculateSaleSplit(input);
    if (split.merchantNetUsd <= 0) {
      return {
        ok: false,
        reason: "El margen neto del comercio es cero o negativo. Ajusta el precio de venta.",
      };
    }
    return { ok: true };
  } catch (error) {
    return {
      ok: false,
      reason: error instanceof Error ? error.message : "Entrada inválida",
    };
  }
}

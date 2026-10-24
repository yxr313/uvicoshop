export function formatUsd(amount: number): string {
  const formatted = new Intl.NumberFormat("es-VE", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
  return `$${formatted}`;
}

export function formatVes(amountUsd: number, bcvRateVesPerUsd: number): string {
  if (!Number.isFinite(bcvRateVesPerUsd) || bcvRateVesPerUsd <= 0) {
    throw new Error("bcvRateVesPerUsd must be a positive finite number");
  }
  const amountVes = amountUsd * bcvRateVesPerUsd;
  const formatted = new Intl.NumberFormat("es-VE", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amountVes);
  return `Bs. ${formatted}`;
}

export function formatPercent(rate: number): string {
  if (!Number.isFinite(rate) || rate < 0 || rate > 1) {
    throw new Error("rate must be between 0 and 1");
  }
  const percent = rate * 100;
  const formatted = new Intl.NumberFormat("es-VE", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(percent);
  return `${formatted}%`;
}

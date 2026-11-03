import { useQuery } from "@tanstack/react-query";
import { fetchFromApi } from "@/lib/api-client";

interface PaymentMethodStat {
  method: string;
  count: number;
}

interface TopProductStat {
  name: string;
  quantity: number;
  revenue_cents: number;
}

interface AnalyticsSummaryResponse {
  summary: {
    new_customers: number;
    new_customers_trend?: number;
    returning_customers: number;
    returning_customers_trend?: number;
    repeat_rate: number;
    repeat_rate_trend?: number;
    top_payment_method: string;
    top_payment_count: number;
    payment_methods: PaymentMethodStat[];
  };
  top_products: TopProductStat[];
}

interface TrafficResponse {
  sources: { source: string; value: number }[];
}

interface VisitsResponse {
  visits: { time: string; visits: number }[];
}

export interface DonutData {
  source: string;
  value: number;
  fill: string;
}

// Paleta de colores estricta
const PAYMENT_COLORS: Record<string, string> = {
  pago_movil: "#A7C957", // Verde Salvia
  zelle: "#171717",      // Negro
  usdt: "#525252",       // Gris oscuro
  transferencia: "#a3a3a3",
  efectivo: "#e5e5e5",
};

const TRAFFIC_COLORS = ["#A7C957", "#171717", "#525252", "#a3a3a3", "#e5e5e5"] as const;

const formatMethod = (method: string): string => {
  const map: Record<string, string> = {
    pago_movil: "Pago Móvil",
    zelle: "Zelle",
    usdt: "USDT",
    transferencia: "Transferencia",
    efectivo: "Efectivo",
  };
  return map[method] ?? method ?? "Desconocido";
};

export function useAnalyticsData(range: string = "30d") {
  const summaryQuery = useQuery({
    queryKey: ["analytics-summary"],
    queryFn: () => fetchFromApi<AnalyticsSummaryResponse>("/analytics/summary"),
  });

  const trafficQuery = useQuery({
    queryKey: ["analytics-traffic"],
    queryFn: () => fetchFromApi<TrafficResponse>("/analytics/traffic"),
  });

  const visitsQuery = useQuery({
    queryKey: ["analytics-visits", range],
    queryFn: () => fetchFromApi<VisitsResponse>(`/analytics/visits?range=${range}`),
  });

  const summary = summaryQuery.data?.summary;
  const topProducts = summaryQuery.data?.top_products ?? [];
  const visitsData = visitsQuery.data?.visits ?? [];

  // Mapeo de KPIs con tendencias reales
  const kpis = [
    { 
      label: "Tasa de Recompra", 
      value: `${(summary?.repeat_rate ?? 0).toFixed(1)}%`, 
      change: summary?.repeat_rate_trend ?? 0, 
      trend: (summary?.repeat_rate_trend ?? 0) >= 0 ? ("up" as const) : ("down" as const), 
      hint: "vs. 30 días anteriores" 
    },
    { 
      label: "Clientes Nuevos", 
      value: (summary?.new_customers ?? 0).toString(), 
      change: summary?.new_customers_trend ?? 0, 
      trend: (summary?.new_customers_trend ?? 0) >= 0 ? ("up" as const) : ("down" as const), 
      hint: "vs. 30 días anteriores" 
    },
    { 
      label: "Clientes Recurrentes", 
      value: (summary?.returning_customers ?? 0).toString(), 
      change: summary?.returning_customers_trend ?? 0, 
      trend: (summary?.returning_customers_trend ?? 0) >= 0 ? ("up" as const) : ("down" as const), 
      hint: "vs. 30 días anteriores" 
    },
    { 
      label: "Método Favorito", 
      value: formatMethod(summary?.top_payment_method ?? "-"), 
      change: 0, 
      trend: "up" as const, 
      hint: `${summary?.top_payment_count ?? 0} transacciones` 
    },
  ];

  const paymentData: DonutData[] = (summary?.payment_methods ?? []).map((item) => ({
    source: formatMethod(item.method),
    value: item.count,
    fill: PAYMENT_COLORS[item.method] ?? "#d4d4d4", 
  }));

  const trafficData: DonutData[] = (trafficQuery.data?.sources ?? []).map((item, index) => ({
    source: item.source,
    value: item.value,
    fill: TRAFFIC_COLORS[index % TRAFFIC_COLORS.length] ?? "#e5e5e5",
  }));

  return {
    isLoading: summaryQuery.isLoading || trafficQuery.isLoading,
    isVisitsLoading: visitsQuery.isLoading,
    kpis,
    topProducts,
    paymentData,
    trafficData,
    visitsData,
  };
}
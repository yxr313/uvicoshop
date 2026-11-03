import { useQuery } from "@tanstack/react-query";
import { fetchFromApi } from "@/lib/api-client";
import { createClient } from "@/lib/supabase/client";

export interface OrderSummary {
  id: string;
  customer_id: string;
  customer_name?: string;
  status: string;
  payment_method?: string;
  total_cents: number;
  merchant_net_cents: number;
  currency: string;
  created_at: string;
}

interface OrdersResponse { orders: OrderSummary[]; }
export interface TrafficSource { source: string; value: number; fill: string; }

const MONTH_NAMES = ["Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"] as const;
const getMonthKey = (d: Date) => MONTH_NAMES[d.getMonth()] ?? "Ene";
const calcChange = (current: number, previous: number) => {
  if (previous === 0) return current > 0 ? 100 : 0;
  return Number((((current - previous) / previous) * 100).toFixed(1));
};

export function useDashboardData() {
  const supabase = createClient();

  const query = useQuery({
    queryKey: ["dashboard-orders"],
    queryFn: async () => {
      try {
        return await fetchFromApi<OrdersResponse>("/orders?limit=100");
      } catch (error: any) {
        if (error.code === "missing_tenant") {
          const { data: { user } } = await supabase.auth.getUser();
          const slug = user?.user_metadata?.store_slug;
          const fullName = user?.user_metadata?.full_name || "Mi";
          if (slug) {
            await fetchFromApi("/tenants", {
              method: "POST",
              body: JSON.stringify({ name: `${fullName}'s Store`, slug, plan_id: "founder", interval: "monthly" }),
            });
            return await fetchFromApi<OrdersResponse>("/orders?limit=100");
          }
        }
        throw error;
      }
    },
  });

  const trafficQuery = useQuery({
    queryKey: ["dashboard-traffic"],
    queryFn: () => fetchFromApi<{ sources: { source: string; value: number }[] }>("/analytics/traffic"),
  });

  const orders = query.data?.orders || [];
  const paidOrders = orders.filter((o) => o.status === "paid" || o.status === "fulfilled");

  const now = new Date();
  const startOfCurrentMonth = new Date(now.getFullYear(), now.getMonth(), 1).getTime();
  const startOfPreviousMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1).getTime();
  const endOfPreviousMonth = new Date(now.getFullYear(), now.getMonth(), 0, 23, 59, 59).getTime();

  const currentMonthOrders = paidOrders.filter((o) => new Date(o.created_at).getTime() >= startOfCurrentMonth);
  const previousMonthOrders = paidOrders.filter((o) => {
    const time = new Date(o.created_at).getTime();
    return time >= startOfPreviousMonth && time <= endOfPreviousMonth;
  });

  const currentRevenue = currentMonthOrders.reduce((sum, o) => sum + o.total_cents, 0);
  const prevRevenue = previousMonthOrders.reduce((sum, o) => sum + o.total_cents, 0);
  const currentNet = currentMonthOrders.reduce((sum, o) => sum + o.merchant_net_cents, 0);
  const prevNet = previousMonthOrders.reduce((sum, o) => sum + o.merchant_net_cents, 0);
  const currentOrdersCount = currentMonthOrders.length;
  const prevOrdersCount = previousMonthOrders.length;
  const currentAvgTicket = currentOrdersCount > 0 ? currentRevenue / currentOrdersCount : 0;
  const prevAvgTicket = prevOrdersCount > 0 ? prevRevenue / prevOrdersCount : 0;

  const metrics = [
    { label: "Ventas Totales", value: `$${(currentRevenue / 100).toLocaleString()}`, change: calcChange(currentRevenue, prevRevenue), trend: currentRevenue >= prevRevenue ? "up" : "down", hint: "vs. mes anterior" },
    { label: "Ganancia Neta", value: `$${(currentNet / 100).toLocaleString()}`, change: calcChange(currentNet, prevNet), trend: currentNet >= prevNet ? "up" : "down", hint: "vs. mes anterior" },
    { label: "Pedidos", value: currentOrdersCount.toString(), change: calcChange(currentOrdersCount, prevOrdersCount), trend: currentOrdersCount >= prevOrdersCount ? "up" : "down", hint: "vs. mes anterior" },
    { label: "Ticket Promedio", value: `$${(currentAvgTicket / 100).toFixed(2)}`, change: calcChange(currentAvgTicket, prevAvgTicket), trend: currentAvgTicket >= prevAvgTicket ? "up" : "down", hint: "vs. mes anterior" },
  ] as const;

  const revenueByMonth: Record<string, { revenue: number; profit: number }> = {};
  for (let i = 5; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
    revenueByMonth[getMonthKey(d)] = { revenue: 0, profit: 0 };
  }

  paidOrders.forEach((o) => {
    const date = new Date(o.created_at);
    const entry = revenueByMonth[getMonthKey(date)];
    if (entry) {
      entry.revenue += o.total_cents / 100;
      entry.profit += o.merchant_net_cents / 100;
    }
  });

  const chartData = Object.entries(revenueByMonth).map(([month, data]) => ({ month, ...data }));
  const TRAFFIC_COLORS = ["#A7C957", "#171717", "#525252", "#a3a3a3", "#e5e5e5"];
  const trafficData: TrafficSource[] = (trafficQuery.data?.sources || []).map((item, index) => ({
    ...item,
    fill: TRAFFIC_COLORS[index % TRAFFIC_COLORS.length] ?? "#e5e5e5",
  }));

  return {
    ...query,
    isLoading: query.isLoading || trafficQuery.isLoading,
    metrics,
    recentOrders: orders.slice(0, 5),
    chartData,
    trafficData,
  };
}
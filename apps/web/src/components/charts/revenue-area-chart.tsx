"use client";

import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

const config = {
  revenue: { label: "Ingresos", color: "#A7C957" },
  profit: { label: "Ganancia", color: "#171717" },
} satisfies ChartConfig;

export function RevenueAreaChart({ data = [] }: { data?: any[] }) {
  return (
    <ChartContainer config={config} className="aspect-auto h-[280px] w-full">
      <AreaChart data={data} margin={{ left: 4, right: 8, top: 8 }}>
        <defs>
          <linearGradient id="fillRevenue" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#A7C957" stopOpacity={0.3} />
            <stop offset="100%" stopColor="#A7C957" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="fillProfit" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#171717" stopOpacity={0.2} />
            <stop offset="100%" stopColor="#171717" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} strokeDasharray="3 3" stroke="#e5e5e5" />
        <XAxis dataKey="month" tickLine={false} axisLine={false} tickMargin={10} tick={{ fill: "#737373", fontSize: 12 }} />
        <YAxis tickLine={false} axisLine={false} tickMargin={8} width={44} tickFormatter={(v) => `$${v}k`} tick={{ fill: "#737373", fontSize: 12 }} />
        <ChartTooltip cursor={false} content={<ChartTooltipContent indicator="dot" />} />
        <Area dataKey="profit" type="monotone" fill="url(#fillProfit)" stroke="#171717" strokeWidth={2} stackId="a" />
        <Area dataKey="revenue" type="monotone" fill="url(#fillRevenue)" stroke="#A7C957" strokeWidth={2} stackId="b" />
      </AreaChart>
    </ChartContainer>
  );
}
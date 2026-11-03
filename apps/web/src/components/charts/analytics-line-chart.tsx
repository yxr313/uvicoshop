"use client";

import { Line, LineChart, CartesianGrid, XAxis, YAxis } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

const config = {
  visits: { label: "Visitantes", color: "#A7C957" },
} satisfies ChartConfig;

export function AnalyticsLineChart({ data = [] }: { data?: any[] }) {
  if (data.length === 0) {
    return (
      <div className="flex h-[300px] items-center justify-center text-sm text-muted-foreground">
        No hay visitas registradas en este período.
      </div>
    );
  }

  return (
    <ChartContainer config={config} className="aspect-auto h-[300px] w-full">
      <LineChart data={data} margin={{ left: 4, right: 8, top: 8 }}>
        <CartesianGrid vertical={false} strokeDasharray="3 3" stroke="#e5e5e5" />
        <XAxis 
          dataKey="time" 
          tickLine={false} 
          axisLine={false} 
          tickMargin={10} 
          tick={{ fill: "#737373", fontSize: 12 }}
        />
        <YAxis 
          tickLine={false} 
          axisLine={false} 
          tickMargin={8} 
          width={44} 
          tick={{ fill: "#737373", fontSize: 12 }}
        />
        <ChartTooltip cursor={false} content={<ChartTooltipContent indicator="line" />} />
        <Line 
          dataKey="visits" 
          type="monotone" 
          stroke="#A7C957" 
          strokeWidth={2} 
          dot={false} 
        />
      </LineChart>
    </ChartContainer>
  );
}
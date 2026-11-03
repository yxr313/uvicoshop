"use client";

import * as React from "react";
import { Label, Pie, PieChart } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

const config = {
  value: { label: "Visitantes" },
} satisfies ChartConfig;

export interface TrafficData {
  source: string;
  value: number;
  fill: string;
}

export function TrafficDonutChart({ data = [] }: { data?: TrafficData[] }) {
  const total = React.useMemo(
    () => data.reduce((a, b) => a + b.value, 0),
    [data]
  );

  if (data.length === 0) {
    return (
      <div className="flex h-[240px] items-center justify-center text-center text-sm text-muted-foreground">
        No hay datos de tráfico disponibles.
      </div>
    );
  }

  return (
    <ChartContainer config={config} className="mx-auto aspect-square h-[240px]">
      <PieChart>
        <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
        <Pie data={data} dataKey="value" nameKey="source" innerRadius={60} strokeWidth={4} stroke="#ffffff">
          <Label
            content={({ viewBox }) => {
              if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                return (
                  <text x={viewBox.cx} y={viewBox.cy} textAnchor="middle">
                    <tspan x={viewBox.cx} y={viewBox.cy} className="fill-foreground text-2xl font-semibold">
                      {total}
                    </tspan>
                    <tspan x={viewBox.cx} y={(viewBox.cy || 0) + 20} className="fill-muted-foreground text-xs">
                      Visitas
                    </tspan>
                  </text>
                );
              }
            }}
          />
        </Pie>
      </PieChart>
    </ChartContainer>
  );
}
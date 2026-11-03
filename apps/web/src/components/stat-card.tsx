import { IconTrendingUp, IconTrendingDown, IconMinus } from "@tabler/icons-react";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

// Definimos la interfaz localmente para no depender de @/data
export interface StatCardData {
  label: string;
  value: string;
  change: number;
  trend: "up" | "down";
  hint: string;
}

export function StatCard({ stat }: { stat: StatCardData }) {
  const isPositive = stat.change > 0;
  const isNegative = stat.change < 0;
  const isNeutral = stat.change === 0;

  return (
    <Card>
      <CardContent className="space-y-2.5 pt-6">
        <p className="text-sm font-medium text-muted-foreground">{stat.label}</p>
        <div className="flex items-end justify-between gap-2">
          <span className="text-3xl font-semibold tracking-tight tabular-nums text-foreground">
            {stat.value}
          </span>
          <span
            className={cn(
              "flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium",
              isPositive && "bg-[#F0F5E6] text-[#95B54D]", // Verde Salvia
              isNegative && "bg-destructive/10 text-destructive", // Rojo
              isNeutral && "bg-neutral-100 text-neutral-500" // Gris neutro
            )}
          >
            {isPositive && <IconTrendingUp className="size-3.5" />}
            {isNegative && <IconTrendingDown className="size-3.5" />}
            {isNeutral && <IconMinus className="size-3.5" />}
            
            {/* Mostramos el número absoluto, pero si es negativo ya tiene el signo */}
            {Math.abs(stat.change)}%
          </span>
        </div>
        <p className="text-xs text-muted-foreground">{stat.hint}</p>
      </CardContent>
    </Card>
  );
}
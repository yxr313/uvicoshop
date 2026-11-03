"use client";

import { useState } from "react";
import { PageHeader } from "@/components/page-header";
import { StatCard } from "@/components/stat-card";
import { DonutChart } from "@/components/charts/donut-chart";
import { AnalyticsLineChart } from "@/components/charts/analytics-line-chart";
import { TopProductsTable } from "@/components/app/top-products-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useAnalyticsData } from "@/hooks/use-analytics";
import { generateAnalyticsPDF } from "@/lib/pdf-generator"; // <-- IMPORTADO
import { IconDownload } from "@tabler/icons-react";
import { toast } from "sonner";

const RANGES = [
  { value: "24h", label: "24 horas" },
  { value: "7d", label: "7 días" },
  { value: "30d", label: "30 días" },
  { value: "1y", label: "1 año" },
];

export default function AnalyticsPage() {
  const [range, setRange] = useState("30d");
  const { isLoading, kpis, topProducts, paymentData, trafficData, visitsData, isVisitsLoading } = useAnalyticsData(range);

  function handleExport() {
    toast.promise(
      new Promise((resolve) => {
        generateAnalyticsPDF({
          tenantName: "Mi Tienda Uvico", // Puedes hacer esto dinámico luego
          range: RANGES.find((r) => r.value === range)?.label.toLowerCase() ?? "30 días",
          kpis,
          paymentData,
          topProducts,
          trafficData,
        });
        setTimeout(resolve, 1000);
      }),
      {
        loading: "Generando documento empresarial...",
        success: "Reporte PDF descargado exitosamente",
        error: "Error al generar el reporte",
      }
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Inteligencia de Negocio"
        description="Analiza el comportamiento de tus clientes y optimiza tu estrategia."
      >
        <Button variant="outline" size="sm" onClick={handleExport}>
          <IconDownload className="size-4 mr-2" /> Exportar PDF
        </Button>
      </PageHeader>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="space-y-2.5 pt-6">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-8 w-32" />
              </CardContent>
            </Card>
          ))
        ) : (
          kpis.map((s) => <StatCard key={s.label} stat={s} />)
        )}
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Métodos de Pago</CardTitle>
            <CardDescription>Preferencia de tus clientes al pagar</CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-center h-[240px]">
            {isLoading ? <Skeleton className="h-[240px] w-[240px] rounded-full" /> : <DonutChart data={paymentData} centerLabel="Pagos" />}
          </CardContent>
          <CardFooter className="flex-col items-stretch gap-2">
            {paymentData.map((t) => (
              <div key={t.source} className="flex items-center gap-2 text-sm">
                <span className="size-2.5 rounded-full" style={{ background: t.fill }} />
                <span className="text-muted-foreground">{t.source}</span>
                <span className="ml-auto font-medium tabular-nums">{t.value}</span>
              </div>
            ))}
          </CardFooter>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Fuentes de Tráfico</CardTitle>
            <CardDescription>De dónde vienen tus compradores</CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-center h-[240px]">
             {isLoading ? <Skeleton className="h-[240px] w-[240px] rounded-full" /> : <DonutChart data={trafficData} centerLabel="Visitas" />}
          </CardContent>
          <CardFooter className="flex-col items-stretch gap-2">
            {trafficData.map((t) => (
              <div key={t.source} className="flex items-center gap-2 text-sm">
                <span className="size-2.5 rounded-full" style={{ background: t.fill }} />
                <span className="text-muted-foreground">{t.source}</span>
                <span className="ml-auto font-medium tabular-nums">{t.value}</span>
              </div>
            ))}
          </CardFooter>
        </Card>
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <TopProductsTable products={topProducts} isLoading={isLoading} />
        
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>Visitantes</CardTitle>
            <CardDescription>
              Tráfico de los últimos {RANGES.find((r) => r.value === range)?.label.toLowerCase() ?? "30 días"}
            </CardDescription>
            <div className="flex flex-wrap gap-1 mt-3">
              {RANGES.map((r) => (
                <button
                  key={r.value}
                  onClick={() => setRange(r.value)}
                  className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
                    range === r.value ? "bg-foreground text-background" : "bg-muted text-muted-foreground hover:bg-muted/80"
                  }`}
                >
                  {r.label}
                </button>
              ))}
            </div>
          </CardHeader>
          <CardContent>
             {isVisitsLoading ? (
                <Skeleton className="h-[300px] w-full" />
             ) : (
                <AnalyticsLineChart data={visitsData} />
             )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
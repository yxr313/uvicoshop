"use client";

import Link from "next/link";
import { PageHeader } from "@/components/page-header";
import { StatCard } from "@/components/stat-card";
import { RevenueAreaChart } from "@/components/charts/revenue-area-chart";
import { TrafficDonutChart } from "@/components/charts/traffic-donut-chart";
import { RecentOrdersTable } from "@/components/app/recent-orders-table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useDashboardData } from "@/hooks/use-dashboard";

export default function DashboardPage() {
  const { metrics, recentOrders, chartData, trafficData, isLoading } = useDashboardData();

  return (
    <div className="space-y-6">
      <PageHeader
        title="Resumen de tu tienda"
        description="Esto es lo que está pasando en tu negocio hoy."
      >
        <Button variant="outline" asChild>
          <Link href="/dashboard/ecommerce/products">Ver catálogo</Link>
        </Button>
        <Button asChild>
          <Link href="/dashboard/ecommerce/products">Crear producto</Link>
        </Button>
      </PageHeader>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="space-y-2.5 pt-6">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-8 w-32" />
                <Skeleton className="h-3 w-40" />
              </CardContent>
            </Card>
          ))
        ) : (
          metrics.map((s) => <StatCard key={s.label} stat={s} />)
        )}
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Ventas y Ganancias</CardTitle> 
            <CardDescription>Rendimiento mensual (GMV vs Neto)</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <Skeleton className="h-[280px] w-full" />
            ) : (
              <RevenueAreaChart data={chartData} />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Fuentes de Tráfico</CardTitle>
            <CardDescription>De dónde vienen tus visitantes</CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-center h-[280px]">
             {isLoading ? (
                <Skeleton className="h-[240px] w-[240px] rounded-full" />
             ) : trafficData.length > 0 ? (
                <TrafficDonutChart data={trafficData} />
             ) : (
                <div className="text-center text-muted-foreground text-sm">
                  <p>Aún no hay visitas registradas.</p>
                  <p className="text-xs mt-1">Los datos aparecerán cuando los clientes visiten tu tienda.</p>
                </div>
             )}
          </CardContent>
        </Card>
      </div>

      <RecentOrdersTable orders={recentOrders} isLoading={isLoading} />
    </div>
  );
}
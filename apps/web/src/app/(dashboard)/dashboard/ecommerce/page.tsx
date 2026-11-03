"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import {
  IconAlertTriangle,
  IconClock,
  IconPackage,
  IconArrowUpRight,
  IconPlus,
  IconBolt,
} from "@tabler/icons-react";

import { PageHeader } from "@/components/page-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { fetchFromApi } from "@/lib/api-client";

// Tipos actualizados
interface OperationalStats {
  pending_payments: number;
  pending_shipments: number;
  low_stock_products: number;
  has_own_inventory: boolean;
  has_dropshipping: boolean;
}

interface LowStockProduct {
  id: string;
  name: string;
  stock: number;
  sku: string;
}

interface PendingOrder {
  id: string;
  customer_name: string;
  total_cents: number;
  payment_method: string;
  sale_mode: string;
  created_at: string;
}

interface OperationsResponse {
  stats: OperationalStats;
  low_stock: LowStockProduct[];
  pending_orders: PendingOrder[];
}

export default function EcommerceOperationsPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["ecommerce-operations"],
    queryFn: () => fetchFromApi<OperationsResponse>("/analytics/operations"),
  });

  const stats = data?.stats;
  const lowStock = data?.low_stock || [];
  const pendingOrders = data?.pending_orders || [];

  // Determinar cuántas columnas usar
  const showOwnInventory = stats?.has_own_inventory || false;
  const showDropshipping = stats?.has_dropshipping || false;
  const isHybrid = showOwnInventory && showDropshipping;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Centro de Operaciones"
        description="Gestiona el día a día de tu tienda: envíos, pagos y stock."
      >
        <Button variant="outline" size="sm" asChild>
          <Link href="/dashboard/ecommerce/orders">
            Ver todos los pedidos
          </Link>
        </Button>
        <Button size="sm" asChild>
          <Link href="/dashboard/ecommerce/products">
            <IconPlus className="size-4 mr-2" /> Añadir producto
          </Link>
        </Button>
      </PageHeader>

      {/* Tarjetas de Alerta Operativa Dinámicas */}
      <div className={`grid gap-4 sm:grid-cols-2 ${isHybrid ? 'lg:grid-cols-4' : 'lg:grid-cols-3'}`}>
        {isLoading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="h-28 flex items-center justify-center">
                <Skeleton className="h-10 w-full" />
              </CardContent>
            </Card>
          ))
        ) : (
          <>
            {showOwnInventory && (
              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Pagos por Verificar</CardTitle>
                  <IconClock className="size-4 text-amber-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats?.pending_payments || 0}</div>
                  <p className="text-xs text-muted-foreground">Comprobantes manuales</p>
                </CardContent>
              </Card>
            )}

            {showDropshipping && (
              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium">Ventas Automáticas</CardTitle>
                  <IconBolt className="size-4 text-emerald-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-emerald-600">
                    {/* Aquí podrías contar las ventas de hoy de dropshipping, por ahora usamos pending_shipments como proxy */}
                    {stats?.pending_shipments || 0}
                  </div>
                  <p className="text-xs text-muted-foreground">Pagadas vía Pasarela</p>
                </CardContent>
              </Card>
            )}

            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Por Despachar</CardTitle>
                <IconPackage className="size-4 text-blue-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats?.pending_shipments || 0}</div>
                <p className="text-xs text-muted-foreground">Pedidos listos para enviar</p>
              </CardContent>
            </Card>

            {showOwnInventory && (
              <Card>
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-destructive">Stock Crítico</CardTitle>
                  <IconAlertTriangle className="size-4 text-destructive" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stats?.low_stock_products || 0}</div>
                  <p className="text-xs text-muted-foreground">Tu inventario propio</p>
                </CardContent>
              </Card>
            )}
          </>
        )}
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        {/* Tabla de Acción Inmediata: Pedidos Recientes */}
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>Acción Requerida: Pedidos</CardTitle>
            <CardDescription>Los últimos pedidos que necesitan tu atención.</CardDescription>
          </CardHeader>
          <CardContent className="px-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-6">Pedido</TableHead>
                  <TableHead>Método</TableHead>
                  <TableHead className="pr-6 text-right">Monto</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={3} className="h-24 text-center">Cargando...</TableCell>
                  </TableRow>
                ) : pendingOrders.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                      Todo al día. No hay pedidos pendientes.
                    </TableCell>
                  </TableRow>
                ) : (
                  pendingOrders.map((o) => (
                    <TableRow key={o.id}>
                      <TableCell className="pl-6 font-mono text-xs">
                        #{o.id.split("-")[0]}
                      </TableCell>
                      <TableCell className="text-xs">
                        <div className="flex items-center gap-2">
                          <span className="capitalize text-muted-foreground">
                            {o.payment_method.replace("_", " ")}
                          </span>
                          <Badge variant="outline" className={o.sale_mode === "own_inventory" ? "text-[10px] text-amber-600 border-amber-200 bg-amber-50" : "text-[10px] text-emerald-600 border-emerald-200 bg-emerald-50"}>
                            {o.sale_mode === "own_inventory" ? "Manual" : "Auto"}
                          </Badge>
                        </div>
                      </TableCell>
                      <TableCell className="pr-6 text-right font-medium">
                        ${(o.total_cents / 100).toFixed(2)}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
          <CardFooter className="border-t px-6 py-3">
            <Link href="/dashboard/ecommerce/orders" className="text-sm font-medium flex items-center gap-1 hover:underline">
              Ir a gestión de pedidos <IconArrowUpRight className="size-3" />
            </Link>
          </CardFooter>
        </Card>

        {/* Tabla de Inventario Crítico (Solo visible si tiene inventario propio) */}
        {showOwnInventory && (
          <Card className="lg:col-span-1">
            <CardHeader>
              <CardTitle>Inventario Crítico</CardTitle>
              <CardDescription>Productos que se están agotando.</CardDescription>
            </CardHeader>
            <CardContent className="px-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="pl-6">Producto</TableHead>
                    <TableHead>SKU</TableHead>
                    <TableHead className="pr-6 text-right">Stock</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={3} className="h-24 text-center">Cargando...</TableCell>
                    </TableRow>
                  ) : lowStock.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                        Inventario saludable. No hay alertas de stock.
                      </TableCell>
                    </TableRow>
                  ) : (
                    lowStock.map((p) => (
                      <TableRow key={p.id}>
                        <TableCell className="pl-6 font-medium text-sm">
                          {p.name}
                        </TableCell>
                        <TableCell className="text-muted-foreground text-xs">
                          {p.sku}
                        </TableCell>
                        <TableCell className="pr-6 text-right">
                          <Badge variant="outline" className={p.stock === 0 ? "bg-destructive/10 text-destructive" : "bg-amber-500/10 text-amber-600"}>
                            {p.stock} uds
                          </Badge>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </CardContent>
            <CardFooter className="border-t px-6 py-3">
              <Link href="/dashboard/ecommerce/products" className="text-sm font-medium flex items-center gap-1 hover:underline">
                Ir a inventario <IconArrowUpRight className="size-3" />
              </Link>
            </CardFooter>
          </Card>
        )}
      </div>
    </div>
  );
}
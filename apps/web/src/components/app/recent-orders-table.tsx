"use client";

import Link from "next/link";
import { IconArrowUpRight } from "@tabler/icons-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import type { OrderSummary } from "@/hooks/use-dashboard";
import { cn } from "@/lib/utils";

interface RecentOrdersTableProps {
  orders: OrderSummary[];
  isLoading: boolean;
}

export function RecentOrdersTable({ orders, isLoading }: RecentOrdersTableProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Últimos Pedidos</CardTitle>
        <CardDescription>Transacciones recientes de tu tienda</CardDescription>
        <Button variant="ghost" size="sm" asChild className="ml-auto">
          <Link href="/dashboard/ecommerce/orders">
            Ver todos <IconArrowUpRight className="size-4 ml-1" />
          </Link>
        </Button>
      </CardHeader>
      <CardContent className="px-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="pl-6">Pedido</TableHead>
              <TableHead>Cliente</TableHead>
              <TableHead>Método de Pago</TableHead>
              <TableHead>Estado</TableHead>
              <TableHead className="pr-6 text-right">Monto</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: 3 }).map((_, i) => (
                <TableRow key={i}>
                  <TableCell className="pl-6"><Skeleton className="h-6 w-16" /></TableCell>
                  <TableCell><Skeleton className="h-6 w-32" /></TableCell>
                  <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                  <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                  <TableCell className="pr-6 text-right"><Skeleton className="h-6 w-16 ml-auto" /></TableCell>
                </TableRow>
              ))
            ) : orders.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                  No hay pedidos aún. Comparte el enlace de tu tienda para empezar a vender.
                </TableCell>
              </TableRow>
            ) : (
              orders.map((o) => {
                const shortId = o.id ? `#${o.id.split("-")[0]}` : "---";
                // Aquí está el fix: usamos '??' para asegurar que sea un string
                const rawCustomerId = o.customer_id?.split("-")[0] ?? "Cliente";
                const customerName = o.customer_name || rawCustomerId;

                return (
                  <TableRow key={o.id}>
                    <TableCell className="pl-6 font-mono text-xs text-muted-foreground">
                      {shortId}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2.5">
                        <Avatar className="size-7">
                          <AvatarFallback className="bg-neutral-100 text-[10px] font-medium text-neutral-600">
                            {customerName.substring(0, 2).toUpperCase()}
                          </AvatarFallback>
                        </Avatar>
                        <span className="font-medium text-sm">{customerName}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-xs capitalize">
                      {o.payment_method?.replace("_", " ") || "No especificado"}
                    </TableCell>
                    <TableCell>
                      <StatusBadge status={o.status} />
                    </TableCell>
                    <TableCell className="pr-6 text-right font-medium tabular-nums">
                      ${(o.total_cents / 100).toFixed(2)}
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

function StatusBadge({ status }: { status: string }) {
  const map: Record<string, { label: string; className: string }> = {
    paid: { label: "Pagado", className: "bg-[#F0F5E6] text-[#95B54D] border-[#A7C957]/20" },
    fulfilled: { label: "Completado", className: "bg-[#F0F5E6] text-[#95B54D] border-[#A7C957]/20" },
    pending_payment: { label: "Pendiente", className: "bg-neutral-100 text-neutral-600 border-neutral-200" },
    failed: { label: "Fallido", className: "bg-red-50 text-red-600 border-red-200" },
    cancelled: { label: "Cancelado", className: "bg-neutral-50 text-neutral-500 border-neutral-200" },
    refunded: { label: "Reembolsado", className: "bg-neutral-50 text-neutral-500 border-neutral-200" },
  };

  const config = map[status] || {
    label: status,
    className: "bg-neutral-100 text-neutral-600 border-neutral-200",
  };

  return (
    <Badge variant="outline" className={cn("px-2 py-0.5", config.className)}>
      {config.label}
    </Badge>
  );
}
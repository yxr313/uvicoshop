"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

export interface TopProduct {
  name: string;
  quantity: number;
  revenue_cents: number;
}

export function TopProductsTable({ products, isLoading }: { products: TopProduct[], isLoading: boolean }) {
  return (
    <Card className="lg:col-span-2">
      <CardHeader>
        <CardTitle>Top Productos</CardTitle>
        <CardDescription>Los productos que más ingresos te generan.</CardDescription>
      </CardHeader>
      <CardContent className="px-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="pl-6">Producto</TableHead>
              <TableHead className="text-center">Unidades</TableHead>
              <TableHead className="pr-6 text-right">Ingresos</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                  Cargando productos...
                </TableCell>
              </TableRow>
            ) : products.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                  Aún no hay ventas registradas.
                </TableCell>
              </TableRow>
            ) : (
              products.map((p, i) => (
                <TableRow key={p.name}>
                  <TableCell className="pl-6 font-medium">
                    <div className="flex items-center gap-3">
                      <span className="flex size-6 items-center justify-center rounded-full bg-muted text-xs text-muted-foreground">
                        {i + 1}
                      </span>
                      {p.name}
                    </div>
                  </TableCell>
                  <TableCell className="text-center tabular-nums">
                    <Badge variant="secondary">{p.quantity}</Badge>
                  </TableCell>
                  <TableCell className="pr-6 text-right font-medium tabular-nums">
                    ${(p.revenue_cents / 100).toLocaleString()}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
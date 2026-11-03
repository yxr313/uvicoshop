"use client";

import { useState } from "react";
import Link from "next/link"; // <-- Importado para la redirección
import { PageHeader } from "@/components/page-header";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  IconPlus,
  IconSearch,
  IconDotsVertical,
  IconAlertCircle,
  IconBox,
  IconTruckDelivery,
} from "@tabler/icons-react";
import { useProducts, type Product } from "@/hooks/use-products";

// Función para formatear el precio
function formatPrice(cents: number, currency: string) {
  return new Intl.NumberFormat("es-VE", {
    style: "currency",
    currency: currency || "USD",
  }).format(cents / 100);
}

// Helper para el modo de venta
function SaleModeBadge({ mode }: { mode: Product["sale_mode"] }) {
  if (mode === "own_inventory") {
    return (
      <Badge variant="outline" className="bg-neutral-100 text-neutral-600 border-neutral-200 gap-1 font-normal">
        <IconBox className="size-3" /> Propio
      </Badge>
    );
  }
  return (
    <Badge variant="outline" className="bg-[#F0F5E6] text-[#95B54D] border-[#A7C957]/20 gap-1 font-normal">
      <IconTruckDelivery className="size-3" /> Dropshipping
    </Badge>
  );
}

// Helper para el estado del producto
function ProductStatusBadge({ status }: { status: Product["status"] }) {
  const map = {
    active: { label: "Activo", className: "bg-[#F0F5E6] text-[#95B54D] border-[#A7C957]/20" },
    draft: { label: "Borrador", className: "bg-neutral-100 text-neutral-600 border-neutral-200" },
    archived: { label: "Archivado", className: "bg-neutral-50 text-neutral-500 border-neutral-200" },
  };
  const config = map[status] || map.draft;
  
  return (
    <Badge variant="outline" className={`px-2 py-0.5 ${config.className}`}>
      {config.label}
    </Badge>
  );
}

export default function ProductsPage() {
  const [search, setSearch] = useState("");
  const [mode, setMode] = useState("all");
  const [status, setStatus] = useState("all");

  const { data: products = [], isLoading, isError, refetch } = useProducts({ search, mode, status });

  return (
    <div className="space-y-6">
      <PageHeader
        title="Productos"
        description="Gestiona el catálogo real de tu tienda."
      >
        {/* AQUI ESTA LA REDIRECCION */}
        <Button className="bg-accent text-accent-foreground hover:bg-accent-hover" asChild>
          <Link href="/dashboard/ecommerce/products/new">
            <IconPlus className="size-4 mr-2" /> Añadir Producto
          </Link>
        </Button>
      </PageHeader>

      <Card className="rounded-2xl border-border shadow-sm overflow-hidden">
        {/* Barra de Filtros */}
        <div className="p-4 border-b border-border flex flex-col lg:flex-row gap-3 items-center justify-between bg-white">
          <div className="relative w-full max-w-sm">
            <IconSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground size-4" />
            <Input
              type="text"
              placeholder="Buscar por nombre o SKU..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 h-9"
            />
          </div>

          <div className="flex flex-col sm:flex-row gap-2 w-full lg:w-auto">
            {/* Selector de Modo */}
            <Select value={mode} onValueChange={setMode}>
              <SelectTrigger className="w-full sm:w-[180px] h-9">
                <SelectValue placeholder="Modo de venta" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos los modos</SelectItem>
                <SelectItem value="own_inventory">Mi Inventario</SelectItem>
                <SelectItem value="dropshipping_verified">Dropshipping</SelectItem>
              </SelectContent>
            </Select>

            {/* Selector de Estado */}
            <Select value={status} onValueChange={setStatus}>
              <SelectTrigger className="w-full sm:w-[150px] h-9">
                <SelectValue placeholder="Estado" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos los estados</SelectItem>
                <SelectItem value="active">Activos</SelectItem>
                <SelectItem value="draft">Borradores</SelectItem>
                <SelectItem value="archived">Archivados</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-6 space-y-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : isError ? (
            <div className="h-[300px] flex flex-col items-center justify-center text-destructive gap-2">
              <IconAlertCircle className="size-8" />
              <p className="text-sm font-medium">Error al cargar los productos.</p>
              <Button variant="outline" size="sm" onClick={() => refetch()}>Reintentar</Button>
            </div>
          ) : products.length === 0 ? (
            <div className="h-[300px] flex flex-col items-center justify-center text-muted-foreground gap-2">
              <p className="text-sm">No se encontraron productos.</p>
              {(search || mode !== "all" || status !== "all") && (
                <Button variant="link" onClick={() => { setSearch(""); setMode("all"); setStatus("all"); }}>
                  Limpiar filtros
                </Button>
              )}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-6">Producto</TableHead>
                  <TableHead>SKU</TableHead>
                  <TableHead>Modo</TableHead>
                  <TableHead>Stock</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead className="text-right">Precio</TableHead>
                  <TableHead className="pr-6"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {products.map((p) => (
                  <TableRow key={p.id}>
                    <TableCell className="pl-6 py-3">
                      <div className="flex items-center gap-3">
                        <div className="size-10 rounded-lg bg-neutral-100 flex items-center justify-center text-xs text-muted-foreground overflow-hidden">
                          {p.name.substring(0, 2).toUpperCase()}
                        </div>
                        <div className="min-w-0">
                          <p className="font-medium text-foreground text-sm truncate">{p.name}</p>
                          <p className="text-xs text-muted-foreground line-clamp-1 max-w-[200px]">
                            {p.description || "Sin descripción"}
                          </p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-xs font-mono">
                      {p.sku}
                    </TableCell>
                    <TableCell>
                      <SaleModeBadge mode={p.sale_mode} />
                    </TableCell>
                    <TableCell>
                      {p.sale_mode === "own_inventory" ? (
                        <span className={p.stock <= 5 ? "text-destructive font-medium text-sm" : "text-foreground text-sm"}>
                          {p.stock} uds
                        </span>
                      ) : (
                        <span className="text-muted-foreground text-xs italic">Auto</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <ProductStatusBadge status={p.status} />
                    </TableCell>
                    <TableCell className="text-right font-medium text-sm">
                      {formatPrice(p.price_cents, p.currency)}
                    </TableCell>
                    <TableCell className="pr-6 text-right">
                      <Button variant="ghost" size="icon" className="size-8">
                        <IconDotsVertical className="size-4 text-muted-foreground" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
"use client";

import { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent } from "@/components/ui/card";
import { 
  IconShieldCheck, IconStar, IconMapPin, IconLink, 
  IconShoppingBag, IconChartBar, IconFileText, IconArrowLeft,
  IconTruckDelivery, IconClock, IconTrendingUp
} from "@tabler/icons-react";
import { toast } from "sonner";
import type { Supplier } from "@/components/app/suppliers/types";

// Simulamos la data
const getMockSupplier = (id: string): Supplier => {
  return {
    id: id,
    name: id === "textiles-maracay" ? "Textiles Maracay" : "Calzados Caracas C.A.",
    description: "Fabricantes de ropa de alta calidad. Envíos a todo el país con los mejores tiempos de entrega.",
    category: "Ropa",
    state: "Carabobo",
    address: "Zona Industrial San Vicente, Galpón 4",
    rif: "J-12345678-9",
    isVerified: true,
    trustScore: 4.8,
    slaHours: 24,
    logoInitials: "TM",
    connectionsCount: 145,
    productsCount: 24,
    completedOrders: 1250,
    responseTimeHours: 2,
    returnRate: 1.5,
    products: [
      { id: "camiseta-oversize", name: "Camiseta Oversize", description: "Algodón peinado 240gsm", priceCents: 800, rating: 4.9, monthlyRotation: 500, marketAdoption: 85 },
      { id: "franela-cuello", name: "Franela Cuello Redondo", description: "Corte clásico, ideal para estampados", priceCents: 650, rating: 4.7, monthlyRotation: 320, marketAdoption: 72 }
    ],
    documents: [
      { id: "d1", title: "Registro Mercantil (Actualizado)", type: "Registro Mercantil", url: "#" },
      { id: "d2", title: "Políticas de Devolución v2", type: "Políticas", url: "#" }
    ]
  };
};

export default function SupplierProfilePage({ params }: { params: { supplierId: string } }) {
  const supplier = getMockSupplier(params.supplierId);

  const handleConnect = () => {
    toast.success(`Solicitud de Match enviada a ${supplier.name}`, {
      description: "Revisa tu bandeja de 'Muestras' para pagar el pedido de prueba.",
    });
  };

  return (
    <div className="max-w-6xl mx-auto pb-12 space-y-6">
      
      {/* Botón Volver */}
      <div>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground -ml-2 hover:bg-neutral-100 hover:text-foreground">
          <Link href="/dashboard/suppliers/dropshipping">
            <IconArrowLeft className="size-4 mr-2" /> Volver al Directorio
          </Link>
        </Button>
      </div>

      {/* HEADER B2B (Estilo Shopify + Uvico) */}
      <div className="bg-white rounded-2xl border border-border shadow-sm overflow-hidden">
        <div className="p-6 md:p-8">
          <div className="flex flex-col md:flex-row md:items-start justify-between gap-6">
            
            {/* Info Principal */}
            <div className="flex gap-4 items-start">
              <div className="size-16 rounded-xl bg-neutral-100 flex items-center justify-center text-2xl font-bold text-neutral-600 border border-border shrink-0">
                {supplier.logoInitials}
              </div>
              <div>
                <div className="flex items-center gap-2 mb-1">
                  <h1 className="text-2xl font-display font-semibold text-foreground">{supplier.name}</h1>
                  {supplier.isVerified && (
                    <Badge className="bg-[#F0F5E6] text-[#95B54D] hover:bg-[#F0F5E6] border-[#A7C957]/30 gap-1 rounded-md px-2 py-0.5">
                      <IconShieldCheck className="size-3.5" /> Verificado
                    </Badge>
                  )}
                </div>
                <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground mt-2">
                  <span className="flex items-center gap-1 bg-neutral-100 px-2 py-1 rounded text-xs font-medium text-foreground">
                    RIF: {supplier.rif}
                  </span>
                  <span className="flex items-center gap-1">
                    <IconMapPin className="size-4" /> {supplier.state}, Venezuela
                  </span>
                  <span className="flex items-center gap-1">
                    <IconStar className="size-4 text-[#A7C957]" /> {supplier.trustScore.toFixed(1)} Trust Score
                  </span>
                </div>
                <p className="text-sm text-muted-foreground mt-3 max-w-2xl leading-relaxed">
                  {supplier.description}
                </p>
              </div>
            </div>

            {/* Botón de Acción */}
            <div className="shrink-0">
              <Button 
                onClick={handleConnect} 
                className="w-full md:w-auto bg-accent text-accent-foreground hover:bg-accent-hover min-w-[160px] h-11 font-medium"
              >
                <IconLink className="size-4 mr-2" /> Hacer Match
              </Button>
            </div>
          </div>

          {/* Barra de Métricas (SLA, Pedidos, etc.) */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8 pt-6 border-t border-border">
            <div>
              <p className="text-xs text-muted-foreground mb-1 flex items-center gap-1"><IconTruckDelivery className="size-3.5" /> Despacho</p>
              <p className="font-medium text-foreground">{supplier.slaHours} Horas</p>
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-1 flex items-center gap-1"><IconShoppingBag className="size-3.5" /> Pedidos Completados</p>
              <p className="font-medium text-foreground">{supplier.completedOrders}</p>
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-1 flex items-center gap-1"><IconClock className="size-3.5" /> Tiempo de Respuesta</p>
              <p className="font-medium text-foreground">{supplier.responseTimeHours} Horas</p>
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-1 flex items-center gap-1"><IconTrendingUp className="size-3.5" /> Tasa de Devolución</p>
              <p className="font-medium text-foreground">{supplier.returnRate}%</p>
            </div>
          </div>
        </div>
      </div>

      {/* TABS INTERNAS (Catálogo / Reputación / Legal) */}
      <Tabs defaultValue="products" className="w-full">
        <div className="border-b border-border mb-6">
          <TabsList className="w-full justify-start h-auto p-0 bg-transparent gap-6">
            <TabsTrigger 
              value="products" 
              className="rounded-none border-b-2 border-transparent data-[state=active]:border-foreground data-[state=active]:bg-transparent px-0 pb-3 text-sm font-medium text-muted-foreground data-[state=active]:text-foreground transition-none shadow-none"
            >
              <IconShoppingBag className="size-4 mr-2" /> Catálogo
            </TabsTrigger>
            <TabsTrigger 
              value="reputation" 
              className="rounded-none border-b-2 border-transparent data-[state=active]:border-foreground data-[state=active]:bg-transparent px-0 pb-3 text-sm font-medium text-muted-foreground data-[state=active]:text-foreground transition-none shadow-none"
            >
              <IconChartBar className="size-4 mr-2" /> Reputación
            </TabsTrigger>
            <TabsTrigger 
              value="documents" 
              className="rounded-none border-b-2 border-transparent data-[state=active]:border-foreground data-[state=active]:bg-transparent px-0 pb-3 text-sm font-medium text-muted-foreground data-[state=active]:text-foreground transition-none shadow-none"
            >
              <IconFileText className="size-4 mr-2" /> Legal y Documentación
            </TabsTrigger>
          </TabsList>
        </div>

        {/* PESTAÑA 1: PRODUCTOS */}
        <TabsContent value="products" className="pt-2">
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {supplier.products.map((product) => (
              <Link key={product.id} href={`/dashboard/suppliers/dropshipping/${supplier.id}/products/${product.id}`}>
                <Card className="h-full hover:border-foreground/30 hover:shadow-sm transition-all cursor-pointer group bg-white">
                  <CardContent className="p-0">
                    <div className="aspect-square bg-neutral-50 flex items-center justify-center relative border-b border-border/50">
                      <IconShoppingBag className="size-12 text-neutral-300 group-hover:scale-105 transition-transform" />
                      <Badge variant="outline" className="absolute top-3 right-3 bg-white text-foreground font-mono border-border">
                        ${(product.priceCents / 100).toFixed(2)}
                      </Badge>
                    </div>
                    <div className="p-4">
                      <h3 className="font-semibold text-foreground truncate">{product.name}</h3>
                      <div className="flex items-center justify-between mt-2 text-xs text-muted-foreground">
                        <span className="flex items-center gap-1"><IconStar className="size-3 text-[#A7C957] fill-current" /> {product.rating}</span>
                        <span>{product.monthlyRotation} rotación</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </TabsContent>

        {/* PESTAÑA 2: REPUTACIÓN */}
        <TabsContent value="reputation" className="pt-2">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <Card className="bg-white border-border">
              <CardContent className="p-6">
                <p className="text-sm font-medium text-muted-foreground mb-1">Total Ventas Mensuales</p>
                <p className="text-3xl font-display font-semibold text-foreground">{supplier.completedOrders}</p>
              </CardContent>
            </Card>
            <Card className="bg-white border-border">
              <CardContent className="p-6">
                <p className="text-sm font-medium text-muted-foreground mb-1">Nivel de Adopción</p>
                <p className="text-3xl font-display font-semibold text-foreground">Alto</p>
                <p className="text-xs text-muted-foreground mt-1">Basado en ventas históricas</p>
              </CardContent>
            </Card>
            <Card className="bg-white border-border">
              <CardContent className="p-6">
                <p className="text-sm font-medium text-muted-foreground mb-1">Conexiones Activas</p>
                <p className="text-3xl font-display font-semibold text-foreground">{supplier.connectionsCount}</p>
              </CardContent>
            </Card>
          </div>
          <Card className="bg-white border-border">
            <CardContent className="p-6 space-y-4">
              <div>
                <h4 className="font-medium text-foreground mb-1">Dirección Física</h4>
                <p className="text-sm text-muted-foreground">{supplier.address}, {supplier.state}</p>
              </div>
              <div className="pt-4 border-t border-border">
                <h4 className="font-medium text-foreground mb-1">Política de Devoluciones</h4>
                <p className="text-sm text-muted-foreground">El proveedor acepta devoluciones solo por fallas de fabricación comprobadas. Los costos de envío por devolución corren por cuenta de Uvico en caso de ser error del proveedor.</p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* PESTAÑA 3: DOCUMENTOS */}
        <TabsContent value="documents" className="pt-2">
          <div className="space-y-3">
            {supplier.documents.map((doc) => (
              <Card key={doc.id} className="bg-white border-border hover:border-foreground/30 transition-colors">
                <CardContent className="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                  <div className="flex items-center gap-3">
                    <div className="p-2.5 bg-neutral-100 rounded-lg">
                      <IconFileText className="size-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="font-medium text-foreground text-sm">{doc.title}</p>
                      <p className="text-xs text-muted-foreground">{doc.type}</p>
                    </div>
                  </div>
                  <Button variant="outline" size="sm" className="border-border text-foreground hover:bg-neutral-100" onClick={() => toast("Abriendo documento...")}>
                    Ver PDF
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
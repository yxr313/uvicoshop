"use client";

import { PageHeader } from "@/components/page-header";
import { Card } from "@/components/ui/card";

export default function PODOrdersPage() {
  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader
        title="Órdenes POD"
        description="Sigue el estado de producción y envío de tus pedidos personalizados."
      />
      <Card className="p-12 flex flex-col items-center justify-center text-center min-h-[400px] border-dashed">
        <h3 className="text-lg font-medium text-foreground mb-2">Bandeja de Producción</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          Cuando un cliente compre un producto POD, la orden llegará aquí para que el proveedor la produzca y envíe.
        </p>
      </Card>
    </div>
  );
}
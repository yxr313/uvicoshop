"use client";

import { PageHeader } from "@/components/page-header";
import { Card } from "@/components/ui/card";

export default function DropshippingConnectionsPage() {
  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader
        title="Mis Conexiones"
        description="Gestiona los proveedores que ya has vinculado a tu tienda."
      />
      <Card className="p-12 flex flex-col items-center justify-center text-center min-h-[400px] border-dashed">
        <h3 className="text-lg font-medium text-foreground mb-2">Proveedores Vinculados</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          Aquí verás la lista de proveedores con los que ya hiciste "Match" y podrás gestionar sus conexiones.
        </p>
      </Card>
    </div>
  );
}
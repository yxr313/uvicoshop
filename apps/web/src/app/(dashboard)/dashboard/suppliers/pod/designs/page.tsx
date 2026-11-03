"use client";

import { PageHeader } from "@/components/page-header";
import { Card } from "@/components/ui/card";

export default function PODDesignsPage() {
  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader
        title="Mis Diseños"
        description="Galería de diseños aprobados y listos para venderse en tu tienda."
      />
      <Card className="p-12 flex flex-col items-center justify-center text-center min-h-[400px] border-dashed">
        <h3 className="text-lg font-medium text-foreground mb-2">Galería de Diseños</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          Aquí verás tus diseños que han pasado la revisión técnica del proveedor y están listos para ser publicados.
        </p>
      </Card>
    </div>
  );
}
"use client";

import { PageHeader } from "@/components/page-header";
import { Card } from "@/components/ui/card";

export default function DropshippingSamplesPage() {
  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <PageHeader
        title="Pedidos de Prueba"
        description="Haz seguimiento a las muestras solicitadas antes de publicar productos."
      />
      <Card className="p-12 flex flex-col items-center justify-center text-center min-h-[400px] border-dashed">
        <h3 className="text-lg font-medium text-foreground mb-2">Bandeja de Muestras</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          Una vez que el proveedor entregue la muestra, se desbloqueará la creación de productos con su catálogo.
        </p>
      </Card>
    </div>
  );
}